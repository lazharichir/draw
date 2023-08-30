// import * as PIXI from "pixi.js";
import { BaseTexture, settings } from "@pixi/core";
import { SCALE_MODES } from "@pixi/constants";
import { Canvas } from "./canvas/Canvas";
import { Palette } from "./canvas/Palette";
import { ErrorBoundary } from "./ErrorBoundary";
import { useResize } from "./hooks/useResize";
import { useAtom, useAtomValue } from "jotai";
import { eraserSelectedAtom, paletteColorsAtom, selectedColorAtom } from "./stores/jotai";
import { useState } from "react";

// set some global pixi settings
settings.RESOLUTION = window.devicePixelRatio || 1;
settings.ROUND_PIXELS = true;
BaseTexture.defaultOptions.scaleMode = SCALE_MODES.NEAREST;

const parseXYZFromUrl = (url: URL) => {
	const pieces = url.pathname.split("/");
	const cid = parseInt(pieces[2]);
	const z = parseInt(pieces[3]);
	const x = parseInt(pieces[4]);
	const y = parseInt(pieces[5]);
	return {
		CanvasID: cid,
		X: x,
		Y: y,
		Z: z,
	};
};

export const App = () => {
	const urldata = parseXYZFromUrl(new URL(window.location.href));

	const [screenWidth, screenHeight] = useResize();
	const colorChoices = useAtomValue(paletteColorsAtom);
	const [eraserSelected, setEraserSelected] = useAtom(eraserSelectedAtom);
	const [selectedColor, setSelectedColor] = useAtom(selectedColorAtom);
	const [[x, y, z], setXYZ] = useState([urldata.X, urldata.Y, urldata.Z]);

	return (
		<div className="fixed top-0 left-0 w-full h-full">
			<ErrorBoundary>
				<Canvas
					screenWidth={screenWidth}
					screenHeight={screenHeight}
					currentBrushColor={selectedColor}
					eraserSelected={eraserSelected}
					backgroundColor={{
						r: 255,
						g: 240,
						b: 229,
						a: 1,
					}}
					gridColor={{
						r: 223,
						g: 219,
						b: 217,
						a: 1,
					}}
					side={1024}
					initialScale={z}
					initialX={x}
					initialY={y}
					onNewXYZ={(x, y, w) => setXYZ([x, y, w])}
				/>
				<div className="fixed bottom-0 left-0 w-full h-fit shadow z-50">
					<div className="inline-flex flex-row items-center justify-center px-4 py-5 text-center bg-white m-auto">
						<div className="text-xs text-gray-500">{x}</div>,
						<div className="text-xs text-gray-500">{y}</div>,
						<div className="text-xs text-gray-500">{z}</div>
					</div>
				</div>
				<div className="fixed top-0 right-0 bg-white w-fit h-full shadow z-50">
					<Palette
						choices={colorChoices}
						value={selectedColor}
						eraserSelected={eraserSelected}
						onChange={(newColor) => {
							setSelectedColor((state) => (state = newColor));
							setEraserSelected((state) => (state = false));
						}}
						onEraserClick={() => {
							setEraserSelected((state) => (state = !state));
						}}
					/>
				</div>
			</ErrorBoundary>
		</div>
	);
};

export default App;
