import * as PIXI from "pixi.js";
import { Stage, Sprite, Container } from "@pixi/react";
import { Viewport } from "pixi-viewport";
import { ViewportClickedEvent, Viewport as ViewportEl } from "./Viewport";
import { useEffect, useMemo, useRef, useState } from "react";
import { ErasedPixel, PixelDrawn, PixelSnapshot, RGBA, TileData } from "../types";
import { retile } from "./retile";
import { useTilesReducer } from "./useTilesReducer";
import { useEnhancedReducer } from "../hooks/useEnhancedReducer";
import { EventSystem } from "@pixi/events";
import isEqual from "lodash.isequal";

export interface CanvasProps {
	screenWidth: number;
	screenHeight: number;
	side: number;
	currentBrushColor: RGBA | null;
	backgroundColor: RGBA;
	gridColor: RGBA;
	eraserSelected: boolean;
	initialScale?: number;
}

type Line = {
	x: number;
	y: number;
	w: number;
	h: number;
};

const ClientPixelContainer = (props: { snapshot: PixelSnapshot; backgroundColor: RGBA }) => {
	const { snapshot, backgroundColor } = props;
	// console.log(`ClientPixelContainer`, { snapshot, backgroundColor });

	const pixels: JSX.Element[] = useMemo(() => {
		// console.log(`ClientPixelContainer.useMemo`, { snapshot, backgroundColor });
		const els: JSX.Element[] = [];

		let k = 0;

		Object.keys(snapshot).forEach((x) => {
			Object.keys(snapshot[+x]).forEach((y) => {
				const { color, erased } = snapshot[+x][+y];
				els.push(
					<Sprite
						key={k++}
						texture={PIXI.Texture.WHITE}
						position={[+x, +y]}
						tint={erased ? backgroundColor : color}
						width={1}
						height={1}
					/>
				);
			});
		});

		return els;
	}, [snapshot, backgroundColor]);

	return <Container>{pixels}</Container>;
};

const TileContainer = (props: { side: number; tiles: TileData[] }) => {
	const { side, tiles } = props;
	return (
		<Container>
			{tiles.map((tile: TileData) => (
				<Sprite
					key={`${tile.x}x${tile.y}_${tile.anchor}`}
					image={`http:\/\/localhost:1001/tile/${tile.x}x${tile.y}_${side}.png`}
					{...tile}
				/>
			))}
		</Container>
	);
};

const useSurface = (worldWidth: number, worldHeight: number, brushColor: RGBA) => {
	const [tileState, dispatch, getTileState] = useEnhancedReducer(useTilesReducer, { tiles: [] });
	const [[lastPointerX, lastPointerY], setLastPointer] = useState<number[]>([]);
	const [lastClick, setLastClick] = useState<[number, number] | null>(null);
	const [snapshot, setSnapshot] = useState<PixelSnapshot>({});
	const viewportRef = useRef<Viewport>(null);

	const drawPixel = (x: number, y: number, color: RGBA) => {
		setSnapshot((prev) => ({
			...prev,
			[+x]: {
				...(prev[+x] || {}),
				[+y]: {
					at: Date.now(),
					color: color,
					erased: false,
				},
			},
		}));
	};

	const erasePixel = (x: number, y: number) => {
		setSnapshot((prev) => ({
			...prev,
			[+x]: {
				...(prev[+x] || {}),
				[+y]: {
					at: Date.now(),
					color: brushColor,
					erased: true,
				},
			},
		}));
	};

	return {
		drawPixel,
		erasePixel,
	};
};

export const Canvas = (props: CanvasProps) => {
	// the screen size
	const {
		screenWidth,
		screenHeight,
		backgroundColor,
		side,
		initialScale = 1,
		eraserSelected = false,
		currentBrushColor,
	} = props;

	const selectedColor = JSON.parse(JSON.stringify(currentBrushColor)) as RGBA | null;

	// use a reducer for tiles
	const [tileState, dispatch, getTileState] = useEnhancedReducer(useTilesReducer, { tiles: [] });

	// local state
	const [[lastPointerX, lastPointerY], setLastPointer] = useState<number[]>([]);
	const [lastClick, setLastClick] = useState<[number, number] | null>(null);
	const [[centerX, centerY], setLastCenter] = useState<[number, number]>([0, 0]);
	const [snapshot, setSnapshot] = useState<PixelSnapshot>({});
	const [scale, setScale] = useState<number>(initialScale);

	// the viewport ref
	const viewportRef = useRef<Viewport>(null);

	// the world size
	const [worldWidth, worldHeight] = [1000000, 1000000];

	const startRetiling = (viewport: Viewport) => {
		dispatch({
			type: `ADD`,
			payload: retile(viewport, getTileState().tiles, side),
		});
	};

	useEffect(() => {
		const pieces = [`/canvas`, `/${0}`, `/${scale}`, `/${centerX}`, `/${centerY}`];

		// get current page title
		window.history.replaceState(null, "New Page Title", pieces.join(``));
	}, [centerX, centerY, scale]);

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
		startRetiling(e.viewport);
		e.viewport.moveCenter(0, 0);
	};

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

	const handleClick = (e: ViewportClickedEvent) => {
		setLastClick([Math.floor(e.world.x), Math.floor(e.world.y)]);
	};

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

	const lines = useMemo(() => {
		const lines: Line[] = [];
		const thickness = scale < 36 ? 0.03 : 0.02;

		if (scale < 16) return lines;
		if (!viewportRef.current) return lines;

		const x = Math.floor(viewportRef.current.corner.x);
		const y = Math.floor(viewportRef.current.corner.y);
		const width = viewportRef.current.screenWidth;
		const height = viewportRef.current.screenHeight;

		// horizontal lines
		for (let i = y; i < y + height; i++) {
			lines.push({
				x,
				y: i,
				w: width,
				h: thickness,
			});
		}

		// vertical lines
		for (let i = x; i < x + width; i++) {
			lines.push({
				x: i,
				y,
				w: thickness,
				h: height,
			});
		}

		return lines;
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
					eventMode: `static`,
					resolution: window.devicePixelRatio || 1,
					antialias: true,
					resizeTo: window,
					eventFeatures: {
						click: true,
						globalMove: true,
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
					onClicked={handleClick}
					onMoved={handleMoved}
					onZoomedEnd={handleZoomedEnd}
					onInited={handleInit}
					onPointerMoved={handlePointerMoved}
					clampZoomOptions={{
						minScale: 1,
						maxScale: 70,
					}}
				>
					<Container>
						{/* <TileContainer /> */}
						<TileContainer side={side} tiles={tileState.tiles} />
						<ClientPixelContainer snapshot={snapshot} backgroundColor={backgroundColor} />
					</Container>
					<Container>
						<Sprite
							texture={PIXI.Texture.WHITE}
							width={1}
							height={1}
							cursor="crosshair"
							x={lastPointerX}
							y={lastPointerY}
							tint={selectedColor || undefined}
							visible={selectedColor !== null}
						/>
						<Sprite
							texture={PIXI.Texture.WHITE}
							width={1}
							height={1}
							cursor="crosshair"
							x={lastPointerX}
							y={lastPointerY}
							tint={backgroundColor}
							visible={eraserSelected}
						/>
					</Container>
					{scale > 16 && (
						<Container>
							{lines.map((line, i) => (
								<Sprite
									key={i}
									texture={PIXI.Texture.WHITE}
									position={[line.x, line.y]}
									width={line.w}
									height={line.h}
									tint={`#dfdbd9`}
									alpha={1}
								/>
							))}
						</Container>
					)}
				</ViewportEl>
			</Stage>
		</div>
	);
};
