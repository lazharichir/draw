import { useAtomValue } from "jotai";
import { RefObject, useCallback, useEffect, useRef, useState } from "react";
import { ErrorBoundary } from "../ErrorBoundary";
import { erasePixelRemotely, pollPixelsRemotely, putPixelRemotely } from "../api";
import { Canvas } from "../canvas/Canvas";
import { Palette } from "../canvas/Palette";
import { useResize } from "../hooks/useResize";
import { CanvasState, Point, SetPixelData, useCanvasReducer } from "../stores/canvas.store";
import { useInstruments } from "../stores/instruments.store";
import { paletteColorsAtom } from "../stores/jotai";
import { BaseTexture, settings } from "@pixi/core";
import { SCALE_MODES } from "@pixi/constants";
import { useParams } from "react-router-dom";
import { Viewport } from "pixi-viewport";
import Modal from "../components/Modal";
import { useCopyToClipboard } from "../hooks/useCopyToClipboard";
import { toastFail, toastSuccess } from "../stores/toasts";
import { useSmartInterval } from "../hooks/useSmartInterval";

// set some global pixi settings
settings.RESOLUTION = window.devicePixelRatio || 1;
settings.ROUND_PIXELS = true;
BaseTexture.defaultOptions.scaleMode = SCALE_MODES.NEAREST;

type CanvasSingleURLParams = {
	id: string;
	z: string;
	x: string;
	y: string;
};

function buildNewURL(state: CanvasState): string {
	const pieces = [`/canvas`, `/${state.canvasId}`, `/${state.scale}`, `/${state.center.x}`, `/${state.center.y}`];
	return pieces.join(``);
}

const getDateObjectTenMinutesAgo = () => {
	const date = new Date();
	date.setMinutes(date.getMinutes() - 10);
	return date;
};

const getWorldTopLeftBottomRightFromViewport = (
	viewport: Viewport | undefined | null,
	center: Point,
	screenWidth: number,
	screenHeight: number
): [Point, Point] => {
	let pt = [
		{ x: 0, y: 0 },
		{ x: 0, y: 0 },
	] as [Point, Point];
	if (!viewport) return pt;

	let screenWidthInWorldPixels = Math.floor(viewport.screenWidthInWorldPixels);
	let screenHeightInWorldPixels = Math.floor(viewport.screenHeightInWorldPixels);

	let halfScreenWidthInWorldPixels = Math.floor(screenWidthInWorldPixels / 2);
	let halfScreenHeightInWorldPixels = Math.floor(screenHeightInWorldPixels / 2);

	let centerX = Math.floor(viewport.center.x);
	let centerY = Math.floor(viewport.center.y);

	let topLeft = {
		x: centerX - halfScreenWidthInWorldPixels,
		y: centerY - halfScreenHeightInWorldPixels,
	};

	let bottomRight = {
		x: centerX + halfScreenWidthInWorldPixels,
		y: centerY + halfScreenHeightInWorldPixels,
	};

	// console.log(` `);
	// console.log(`center`, centerX, centerY);
	// console.log(`viewport.screenWidth/Height`, viewport.screenWidth, viewport.screenHeight);
	// console.log(`viewport.screenW/HInWorldPxs`, screenWidthInWorldPixels, screenHeightInWorldPixels);
	// console.log(`top left`, topLeft.x, topLeft.y);
	// console.log(`bottom right`, bottomRight.x, bottomRight.y);
	// console.log(` `);

	return [topLeft, bottomRight];
};

export default function CanvasSingle() {
	const params = useParams<CanvasSingleURLParams>();
	const id = parseInt(params.id || "0");
	const z = parseInt(params.z || "1");
	const x = parseInt(params.x || "0");
	const y = parseInt(params.y || "0");

	const instruments = useInstruments();
	const copyToClipboard = useCopyToClipboard();
	const [screenWidth, screenHeight] = useResize();
	const colorChoices = useAtomValue(paletteColorsAtom);
	const canvasState = useCanvasReducer({
		canvasId: id,
		screenWidth,
		screenHeight,
		scale: z,
		center: {
			x: x,
			y: y,
		},
	});

	const [viewportRef, setViewportRef] = useState<RefObject<Viewport>>();
	const [showGoToModal, setShowGoToModal] = useState(false);
	const [topLeft, setTopLeft] = useState<Point>({ x: 0, y: 0 });
	const [bottomRight, setBottomRight] = useState<Point>({ x: 0, y: 0 });
	const [dragging, setIsDragging] = useState<boolean>(false);

	useEffect(() => {
		const [tl, br] = getWorldTopLeftBottomRightFromViewport(
			viewportRef?.current,
			canvasState.state.center,
			canvasState.state.screenWidth,
			canvasState.state.screenHeight
		);

		setTopLeft(tl);
		setBottomRight(br);
	}, [canvasState.state.center, canvasState.state.screenWidth, canvasState.state.screenHeight, viewportRef?.current]);

	const timerIdRef = useRef<number | null>(null);
	const [isPollingEnabled, setIsPollingEnabled] = useState(true);
	const [pollingFrom, setPollingFrom] = useState<Date>(new Date());
	const smartInterval = useSmartInterval({
		min: 1000,
		max: 10000,
		initialValue: 1000,
		successMultiplier: 0.1,
		errorMultiplier: 2,
		quietMultiplier: 1.1,
		quietMax: 5000,
	});

	useEffect(() => {
		const pollingCallback = async () => {
			try {
				const pixels = await pollPixelsRemotely(canvasState.state.canvasId, topLeft, bottomRight, pollingFrom);
				const pixelsWithIds: SetPixelData[] = pixels.map((p) => ({
					id: crypto.randomUUID(),
					at: Date.now(),
					color: {
						r: p.RGBA.R,
						g: p.RGBA.G,
						b: p.RGBA.B,
						a: p.RGBA.A,
					},
					deleted: false,
					x: p.X,
					y: p.Y,
				}));
				setPollingFrom(new Date());

				console.log(
					`> polled pixels (delay: ${smartInterval.value.toLocaleString()}ms)`,
					pixelsWithIds.length,
					pixelsWithIds || []
				);

				if (pixelsWithIds.length === 0) {
					smartInterval.quiet();
				} else {
					canvasState.setPixels(pixelsWithIds);
					smartInterval.success();
				}
			} catch (error) {
				console.error(`> polling failed`, error);
				smartInterval.error();
			}
		};

		const startPolling = () => {
			timerIdRef.current = setInterval(pollingCallback, smartInterval.value);
		};

		const stopPolling = () => {
			clearInterval(timerIdRef.current || undefined);
		};

		if (isPollingEnabled && !dragging) {
			startPolling();
		} else {
			stopPolling();
		}

		return () => {
			stopPolling();
		};
	}, [
		dragging,
		isPollingEnabled,
		smartInterval.value,
		pollingFrom,
		topLeft,
		bottomRight,
		canvasState.state.canvasId,
	]);

	useEffect(() => canvasState.setScreenSize(screenWidth, screenHeight), [screenWidth, screenHeight]);

	useEffect(() => {
		if (dragging) return;
		window.history.replaceState(null, "", buildNewURL(canvasState.state));
	}, [dragging, canvasState.state.canvasId, canvasState.state.center, canvasState.state.scale]);

	function handleCanvasClick(x: number, y: number) {
		const currentMode = instruments.current().mode;
		switch (currentMode) {
			case "eraser":
				doErasePixel(x, y);
				break;
			case "brush":
				doDrawPixel(x, y);
				break;
			default:
				console.error(`Unknown mode: ${currentMode}`);
				break;
		}
	}

	async function doErasePixel(x: number, y: number) {
		const pxs = [
			{
				id: crypto.randomUUID(),
				x: x,
				y: y,
				deleted: true,
				color: instruments.current().brushColor,
				at: Date.now(),
			},
		];
		canvasState.setPixels(pxs);

		try {
			await erasePixelRemotely(canvasState.state.canvasId, x, y);
			// sendMessage(`[ERASE] (${x},${y})`);
		} catch (error) {
			canvasState.unsetPixels(pxs);
			toastFail(`Failed to erase pixel (removal undone from your canvas).`);
		}
	}

	async function doDrawPixel(x: number, y: number) {
		const pxs = [
			{
				id: crypto.randomUUID(),
				x: x,
				y: y,
				deleted: false,
				color: instruments.current().brushColor,
				at: Date.now(),
			},
		];
		canvasState.setPixels(pxs);
		try {
			// sendMessage(`[PAINT] (${x},${y}) (${JSON.stringify(instruments.current().brushColor)}))`);
			await putPixelRemotely(canvasState.state.canvasId, x, y, instruments.current().brushColor);
		} catch (error) {
			canvasState.unsetPixels(pxs);
			toastFail(`Failed to draw pixel (now undone from your canvas).`);
		}
	}

	function setViewportCenter(point: Point) {
		if (!viewportRef) return;
		if (!viewportRef.current) return;
		if (isNaN(point.x) || isNaN(point.y)) return;
		viewportRef.current.moveCenter(point.x, point.y);
		canvasState.setCenter(point.x, point.y);
	}

	return (
		<>
			<div className="fixed top-0 left-0 w-full h-full">
				<ErrorBoundary>
					<div className="fixed top-0 left-0 w-full h-full z-0">
						<Canvas
							state={canvasState.state}
							instruments={instruments.state}
							onClick={handleCanvasClick}
							onCenterChange={(x, y) => canvasState.setCenter(x, y)}
							onScaleChange={(w) => canvasState.setScale(w)}
							onViewportChange={(viewport) => canvasState.retile(viewport)}
							onViewportRefInit={(viewport) => setViewportRef(viewport)}
							onDragStart={() => setIsDragging(true)}
							onDragEnd={() => setIsDragging(false)}
						/>
					</div>
					<div className="fixed top-0 left-0 w-fit h-full z-50 bg-white bg-opacity-50 hidden">
						<button
							className="p-4"
							onClick={async () => {
								setIsPollingEnabled(!isPollingEnabled);
							}}
						>
							{isPollingEnabled ? `Stop Polling` : `Start Polling`}
						</button>
					</div>
					<div className="fixed bottom-0 left-0 w-full h-fit z-50">
						<div className="p-4 flex flex-row gap-4 items-center justify-center">
							<div>
								<div className="flex flex-row items-center justify-center px-4 py-3 text-center bg-white m-auto gap-3 rounded-lg shadow-md leading-none">
									<button
										className="w-4 h-4"
										onClick={() => {
											copyToClipboard[1](window.location.href);
											toastSuccess(`Copied to clipboard!`);
										}}
									>
										<img src="/pin.svg" alt="location" className="w-4 h-4" />
									</button>
									<div className="text-xs text-black">
										<span className="text-gray-400">x</span>{" "}
										<strong className="font-bold">{canvasState.state.center.x}</strong>
									</div>
									<div className="text-xs text-black">
										<span className="text-gray-400">y</span>{" "}
										<strong className="font-bold">{canvasState.state.center.y}</strong>
									</div>
									<div className="text-xs text-black">
										<span className="text-gray-400">z</span>{" "}
										<strong className="font-bold">{canvasState.state.scale}</strong>
									</div>
								</div>
							</div>
							<div>
								<div className="flex flex-row items-center justify-center px-3 py-3 text-center bg-white m-auto gap-3 rounded-lg shadow-md">
									<div className="flex flex-row gap-2 leading-none">
										<button
											onClick={() => {
												setViewportCenter({ x: 0, y: 0 });
												toastSuccess(`We're back to the center!`);
											}}
											className="p-0 m-0 "
										>
											<img src="/center.svg" alt="location" className="w-4 h-4" />
										</button>
										<button className="p-0 m-0 " onClick={() => setShowGoToModal(true)}>
											<img src="/goto.svg" alt="location" className="w-4 h-4" />
										</button>
									</div>
								</div>
							</div>
						</div>
					</div>
					<div className="fixed top-0 right-0 bg-white w-fit h-full shadow z-50">
						<Palette
							choices={colorChoices}
							value={instruments.state.brushColor}
							eraserSelected={instruments.state.mode === `eraser`}
							onChange={(newColor) => {
								if (!newColor) return;
								instruments.setBrushColor(newColor);
								instruments.setMode(`brush`);
							}}
							onEraserClick={() => {
								instruments.setMode(`eraser`);
							}}
						/>
					</div>
					<Modal
						isOpen={showGoToModal}
						onClose={() => setShowGoToModal(false)}
						onPointChange={(newPoint) => {
							setViewportCenter(newPoint);
							toastSuccess(`We're now at (${newPoint.x},${newPoint.y})!`);
						}}
					/>
				</ErrorBoundary>
			</div>
		</>
	);
}

// // websockets
// const [socketUrl, setSocketUrl] = useState(`ws://localhost:1001/ws`);
// const [messageHistory, setMessageHistory] = useState<any[]>([]);
// const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);

// useEffect(() => {
// 	if (!lastMessage) return;
// 	setMessageHistory((prev) => [...prev, lastMessage]);
// }, [lastMessage, setMessageHistory]);

// const connectionStatus = {
// 	[ReadyState.CONNECTING]: "Connecting",
// 	[ReadyState.OPEN]: "Open",
// 	[ReadyState.CLOSING]: "Closing",
// 	[ReadyState.CLOSED]: "Closed",
// 	[ReadyState.UNINSTANTIATED]: "Uninstantiated",
// }[readyState];

/* <div className="fixed top-0 left-0 w-fit h-full z-50 bg-white bg-opacity-50">
<span>The WebSocket is currently {connectionStatus}</span>
<ul className=" max-w-xl">
	{lastMessage ? <li>Last message: {lastMessage.data}</li> : null}
	{messageHistory.map((message, idx) => (
		<li className="block border-t " key={idx}>
			{message ? message.data : null}
		</li>
	))}
</ul>
</div> */
