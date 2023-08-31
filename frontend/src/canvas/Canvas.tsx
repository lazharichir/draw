import { Texture, Point } from "@pixi/core";
import { Stage, Sprite, Container } from "@pixi/react";
import { Viewport } from "pixi-viewport";
import { RefObject, useEffect, useMemo, useRef, useState } from "react";
import { ViewportClickedEvent, Viewport as ViewportEl } from "./Viewport";
import { RGBA } from "../types";
import { CanvasState } from "../stores/canvas.store";
import { InstrumentState } from "../stores/instruments.store";

export interface CanvasProps {
	state: CanvasState;
	instruments: InstrumentState;
	onClick?: (x: number, y: number) => void;
	onCenterChange?: (x: number, y: number) => void;
	onScaleChange?: (scale: number) => void;
	onViewportChange?: (viewport: Viewport) => void;
	onViewportRefInit?: (viewport: RefObject<Viewport>) => void;
}

export const Canvas = (props: CanvasProps) => {
	// the screen size
	const {
		state,
		instruments,
		onClick = () => {},
		onCenterChange = () => {},
		onScaleChange = () => {},
		onViewportChange = () => {},
		onViewportRefInit = () => {},
	} = props;

	const { screenWidth, screenHeight } = state;
	const { side, worldWidth, worldHeight } = state;
	const { brushColor, mode } = instruments;
	const { backgroundColor, gridColor } = state;
	const { center, scale } = state;

	const selectedColor = JSON.parse(JSON.stringify(brushColor)) as RGBA | null;

	// local state
	const [[lastPointerX, lastPointerY], setLastPointer] = useState<number[]>([]);

	// the viewport ref
	const viewportRef = useRef<Viewport>(null);

	useEffect(() => onViewportRefInit(viewportRef), [!!viewportRef.current]);

	useEffect(() => {
		if (!viewportRef.current) return;
		// if (viewportRef.current.x === state.center.x && viewportRef.current.y === state.center.y) return;
		// viewportRef.current.moveCenter(center.x, center.y);
	}, [state.center]);

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

	const handleZoomedEnd = (e: { viewport: Viewport }) => {
		onViewportChange(e.viewport);
		onScaleChange(Math.round(e.viewport.scaled));
	};

	const handleMoved = (e: { viewport: Viewport }) => {
		onViewportChange(e.viewport);
		onCenterChange(Math.round(e.viewport.center.x), Math.round(e.viewport.center.y));
	};

	const handlePointerMoved = (e: Point) => {
		setLastPointer([Math.floor(e.x), Math.floor(e.y)]);
	};

	const handleInit = (e: { viewport: Viewport }) => {
		e.viewport.scale.set(scale);
		e.viewport.moveCenter(center.x, center.y);
		onViewportChange(e.viewport);
	};

	const handleClick = (e: ViewportClickedEvent) => {
		onClick(Math.floor(e.world.x), Math.floor(e.world.y));
	};

	const lines: JSX.Element[] = useMemo(() => {
		const linesEl: JSX.Element[] = [];
		if (scale < 16) return linesEl;
		if (!viewportRef.current) return linesEl;

		const anchor = 0.5;
		const thickness = scale < 40 ? 0.1 : 0.15;
		const x = Math.floor(viewportRef.current.corner.x);
		const y = Math.floor(viewportRef.current.corner.y);
		const width = viewportRef.current.screenWidth;
		const height = viewportRef.current.screenHeight;

		// horizontal lines
		for (let i = y; i < y + height; i++) {
			linesEl.push(
				<Sprite
					cursor="crosshair"
					key={crypto.randomUUID()}
					texture={Texture.WHITE}
					position={[x, i]}
					anchor={anchor}
					width={width}
					height={thickness}
					tint={gridColor}
					alpha={gridColor.a}
				/>
			);
		}

		// vertical lines
		for (let i = x; i < x + width; i++) {
			linesEl.push(
				<Sprite
					cursor="crosshair"
					key={crypto.randomUUID()}
					texture={Texture.WHITE}
					position={[i, y]}
					width={thickness}
					anchor={anchor}
					height={height}
					tint={gridColor}
					alpha={gridColor.a}
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

	const pixels: JSX.Element[] = useMemo(() => {
		const els: JSX.Element[] = [];
		let k = 0;
		Object.keys(state.pixels || {}).forEach((x) => {
			Object.keys(state.pixels[+x] || {}).forEach((y) => {
				const log = state.pixels[+x][+y];
				if (!log.length) return;

				const { color, erased } = log[0];
				els.push(
					<Sprite
						cursor="crosshair"
						key={k++}
						texture={Texture.WHITE}
						x={+x}
						y={+y}
						tint={erased ? backgroundColor : color}
						width={1}
						height={1}
						anchor={0}
					/>
				);
			});
		});

		return els;
	}, [state.pixels]);

	const tiles: JSX.Element[] = useMemo(() => {
		return state.tiles.map((tile) => {
			return (
				<Sprite
					cursor="crosshair"
					key={`${tile.x}x${tile.y}_${tile.anchor}`}
					image={`http:\/\/localhost:1001/tile/${tile.x}x${tile.y}_${side}.png`}
					{...tile}
				/>
			);
		});
	}, [state.tiles]);

	return (
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
					<Container>{tiles}</Container>
					<Container>{pixels}</Container>
				</Container>
				<Container>
					<Sprite
						texture={Texture.WHITE}
						cursor="crosshair"
						width={1}
						height={1}
						x={lastPointerX}
						y={lastPointerY}
						tint={selectedColor || undefined}
						visible={mode === `brush`}
					/>
					<Sprite
						texture={Texture.WHITE}
						cursor="crosshair"
						width={1}
						height={1}
						x={lastPointerX}
						y={lastPointerY}
						tint={backgroundColor}
						visible={mode === `eraser`}
					/>
				</Container>
				<Container>{lines}</Container>
			</ViewportEl>
		</Stage>
	);
};
