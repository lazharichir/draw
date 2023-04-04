import * as PIXI from "pixi.js";
import { EventSystem } from "@pixi/events";
import { Stage, Sprite, Container } from "@pixi/react";
import { Viewport } from "pixi-viewport";
import { useEffect, useMemo, useRef, useState } from "react";
import { ViewportClickedEvent, Viewport as ViewportEl } from "./Viewport";
import { PixelSnapshot, RGBA } from "../types";
import { retile } from "./retile";
import { useTilesReducer } from "./useTilesReducer";
import { useEnhancedReducer } from "../hooks/useEnhancedReducer";
import { ClientPixelContainer } from "./ClientPixelContainer";
import { TileContainer } from "./TileContainer";
import { Grid } from "./Grid";

export interface CanvasProps {
	worldWidth?: number;
	worldHeight?: number;
	screenWidth: number;
	screenHeight: number;
	side: number;
	currentBrushColor: RGBA | null;
	backgroundColor: RGBA;
	gridColor: RGBA;
	eraserSelected: boolean;
	initialScale?: number;
}

const putPixelRemotely = async (canvasID: number, x: number, y: number, pixelColor: RGBA) => {
	try {
		const colr = new PIXI.Color(pixelColor);
		const [r, g, b] = colr.toUint8RgbArray();
		const a = Math.round(colr.alpha * 255);

		const res = await fetch(`http://localhost:1001/pixel/${canvasID}/${x}/${y}/${r}/${g}/${b}/${a}`, {
			method: `PUT`,
			headers: {
				"Content-Type": `application/json`,
			},
		});

		console.log(`putPixelRemotely.res`, res.ok, res);

		if (!res.ok) {
			console.error(`Error putting pixel remotely`, res);
		}
	} catch (error) {
		console.error(`Error putting pixel remotely [catch]`, error);
	}
};

const erasePixelRemotely = async (canvasID: number, x: number, y: number) => {
	try {
		const res = await fetch(`http://localhost:1001/pixel/${canvasID}/${x}/${y}`, {
			method: `DELETE`,
			headers: {
				"Content-Type": `application/json`,
			},
		});

		console.log(`erasePixelRemotely.res`, res.ok, res);

		if (!res.ok) {
			console.error(`Error erasing pixel remotely`, res);
		} else {
			console.log(`pixel erased remotely`);
		}
	} catch (error) {
		console.error(`Error erasing pixel remotely [catch]`, error);
	}
};

export const Canvas = (props: CanvasProps) => {
	// the screen size
	const {
		worldWidth = 9999999999,
		worldHeight = 9999999999,
		screenWidth,
		screenHeight,
		backgroundColor,
		side,
		initialScale = 1,
		eraserSelected = false,
		currentBrushColor,
		gridColor = { r: 0, g: 0, b: 0, a: 0.1 },
	} = props;

	const selectedColor = JSON.parse(JSON.stringify(currentBrushColor)) as RGBA | null;

	// use a reducer for tiles
	const [tileState, dispatch, getTileState] = useEnhancedReducer(useTilesReducer, { tiles: [] });

	// local state
	const [[lastPointerX, lastPointerY], setLastPointer] = useState<number[]>([]);
	const [[centerX, centerY], setLastCenter] = useState<[number, number]>([0, 0]);
	const [lastClick, setLastClick] = useState<[number, number] | null>(null);
	const [snapshot, setSnapshot] = useState<PixelSnapshot>({});
	const [scale, setScale] = useState<number>(initialScale);

	// the viewport ref
	const viewportRef = useRef<Viewport>(null);

	const startRetiling = (viewport: Viewport) => {
		dispatch({
			type: `ADD`,
			payload: retile(viewport, getTileState().tiles, side),
		});
	};

	useEffect(() => {
		const pieces = [`/canvas`, `/${0}`, `/${scale}`, `/${centerX}`, `/${centerY}`];
		window.history.replaceState(null, "New Page Title", pieces.join(``));
	}, [centerX, centerY, scale]);

	useEffect(() => {
		const wheelListener = (e: WheelEvent) => {
			if (e.ctrlKey) e.preventDefault();
		};

		const keydownListener = (e: KeyboardEvent) => {
			if (!e.ctrlKey && !e.metaKey) return;
			if ([`=`, `+`, `_`, `-`, `0`].includes(e.key)) e.preventDefault();
		};

		window.addEventListener("wheel", wheelListener, { passive: false });
		window.addEventListener("keydown", keydownListener, { passive: false });

		return () => {
			window.removeEventListener("wheel", wheelListener);
			window.removeEventListener("keydown", keydownListener);
		};
	}, []);

	useEffect(() => {
		if (!lastClick) return;
		if (!currentBrushColor) return;

		// erase mode
		if (eraserSelected === true) {
			console.log(`eraser selected`, ...lastClick);
			setSnapshot((prev) => ({
				...prev,
				[lastClick[0]]: {
					...(prev[lastClick[0]] || {}),
					[lastClick[1]]: {
						at: Date.now(),
						color: currentBrushColor,
						erased: true,
					},
				},
			}));

			erasePixelRemotely(0, lastClick[0], lastClick[1]);
		}

		// draw mode
		if (eraserSelected === false) {
			setSnapshot((prev) => ({
				...prev,
				[lastClick[0]]: {
					...(prev[lastClick[0]] || {}),
					[lastClick[1]]: {
						at: Date.now(),
						color: currentBrushColor,
						erased: false,
					},
				},
			}));

			putPixelRemotely(0, lastClick[0], lastClick[1], currentBrushColor);
		}
	}, [lastClick]);

	const handleZoomedEnd = (e: { viewport: Viewport }) => {
		startRetiling(e.viewport);
		setScale(Math.round(e.viewport.scaled));
	};

	const handleMoved = (e: { viewport: Viewport }) => {
		startRetiling(e.viewport);
		setLastCenter([Math.round(e.viewport.center.x), Math.round(e.viewport.center.y)]);
	};

	const handlePointerMoved = (e: PIXI.Point) => {
		setLastPointer([Math.floor(e.x), Math.floor(e.y)]);
	};

	const handleInit = (e: { viewport: Viewport }) => {
		e.viewport.moveCenter(0, 0);
		startRetiling(e.viewport);
	};

	const handleClick = (e: ViewportClickedEvent) => {
		setLastClick([Math.floor(e.world.x), Math.floor(e.world.y)]);
	};

	const lines: JSX.Element[] = useMemo(() => {
		const linesEl: JSX.Element[] = [];

		if (scale < 16) return linesEl;
		if (!viewportRef.current) return linesEl;

		const thickness = scale < 36 ? 0.03 : 0.03;
		const x = Math.floor(viewportRef.current.corner.x);
		const y = Math.floor(viewportRef.current.corner.y);
		const width = viewportRef.current.screenWidth;
		const height = viewportRef.current.screenHeight;

		// horizontal lines
		for (let i = y; i < y + height; i++) {
			linesEl.push(
				<Sprite
					key={i}
					texture={PIXI.Texture.WHITE}
					position={[x, i]}
					width={width}
					height={thickness}
					tint={`#dfdbd9`}
				/>
			);
		}

		// vertical lines
		for (let i = x; i < x + width; i++) {
			linesEl.push(
				<Sprite
					key={i}
					texture={PIXI.Texture.WHITE}
					position={[i, y]}
					width={thickness}
					height={height}
					tint={`#dfdbd9`}
				/>
			);
		}

		return linesEl;
	}, [
		Math.floor(viewportRef.current?.corner.x || 0),
		Math.floor(viewportRef.current?.corner.y || 0),
		viewportRef.current?.screenWidth || 0,
		viewportRef.current?.screenHeight || 0,
		scale,
	]);

	return (
		<div className="fixed top-0 left-0 w-full h-full z-0">
			<Stage
				width={screenWidth}
				height={screenHeight}
				options={{
					background: backgroundColor,
					resolution: window.devicePixelRatio || 1,
					resizeTo: window,
					eventMode: `static`,
					antialias: true,
					autoStart: true,
					eventFeatures: {
						globalMove: true,
						click: true,
						move: true,
						wheel: true,
					},
					autoDensity: true,
				}}
			>
				<ViewportEl
					ref={viewportRef}
					worldWidth={worldWidth}
					worldHeight={worldHeight}
					screenWidth={screenWidth}
					screenHeight={screenHeight}
					onInited={handleInit}
					onClicked={handleClick}
					onMoved={handleMoved}
					onZoomedEnd={handleZoomedEnd}
					onPointerMoved={handlePointerMoved}
					clampZoomOptions={{
						minScale: 1,
						maxScale: 70,
					}}
				>
					<Container>
						<TileContainer side={side} tiles={tileState.tiles} />
						<ClientPixelContainer snapshot={snapshot} backgroundColor={backgroundColor} />
					</Container>
					<Container>
						<Sprite
							texture={PIXI.Texture.WHITE}
							cursor="crosshair"
							width={1}
							height={1}
							x={lastPointerX}
							y={lastPointerY}
							tint={selectedColor || undefined}
							visible={selectedColor !== null}
						/>
						<Sprite
							texture={PIXI.Texture.WHITE}
							cursor="crosshair"
							width={1}
							height={1}
							x={lastPointerX}
							y={lastPointerY}
							tint={backgroundColor}
							visible={eraserSelected}
						/>
					</Container>
					{scale >= 16 && <Container>{lines}</Container>}
				</ViewportEl>
				{/* <Container>
					<Grid viewport={viewportRef.current} />
				</Container> */}
			</Stage>
		</div>
	);
};
