// import * as PIXI from "pixi.js";
import { BaseTexture, settings } from "@pixi/core";
import { SCALE_MODES } from "@pixi/constants";
import { Canvas } from "./canvas/Canvas";
import { Palette } from "./canvas/Palette";
import { ErrorBoundary } from "./ErrorBoundary";
import { useResize } from "./hooks/useResize";
import { useAtomValue } from "jotai";
import { paletteColorsAtom } from "./stores/jotai";
import { useEffect, useState } from "react";
import { useCanvasReducer } from "./stores/canvas";

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
	const [[x, y, z], setXYZ] = useState([urldata.X, urldata.Y, urldata.Z]);

	const canvasState = useCanvasReducer({
		canvasId: urldata.CanvasID,
		screenWidth,
		screenHeight,
	});

	useEffect(() => console.log(`canvasState changed: `, canvasState.state), [canvasState.state]);

	useEffect(() => canvasState.setScreenSize(screenWidth, screenHeight), [screenWidth, screenHeight]);

	useEffect(() => {
		const pieces = [
			`/canvas`,
			`/${canvasState.state.canvasId}`,
			`/${canvasState.state.scale}`,
			`/${canvasState.state.center.x}`,
			`/${canvasState.state.center.y}`,
		];
		window.history.replaceState(null, "New Page Title", pieces.join(``));
	}, [canvasState.state.center, canvasState.state.scale]);

	return (
		<div className="fixed top-0 left-0 w-full h-full">
			<ErrorBoundary>
				<div className="fixed top-0 left-0 w-full h-full z-0">
					<Canvas
						state={canvasState.state}
						onClick={(x, y) => canvasState.setLastClick(x, y)}
						onCenterChange={(x, y) => canvasState.setCenter(x, y)}
						onScaleChange={(w) => canvasState.setScale(w)}
					/>
				</div>
				<div className="fixed bottom-0 left-0 w-full h-fit z-50">
					<div className="p-4 flex flex-row gap-4 items-center justify-center">
						<div>
							<div className="flex flex-row items-center justify-center px-4 py-3 text-center bg-white m-auto gap-3 rounded-lg shadow-md">
								<div>
									<img src="/pin.svg" alt="location" className="w-4 h-4" />
								</div>
								<div className="text-xs text-black">
									<span className="text-gray-400">x</span> <strong className="font-bold">{x}</strong>
								</div>
								<div className="text-xs text-black">
									<span className="text-gray-400">y</span> <strong className="font-bold">{y}</strong>
								</div>
								<div className="text-xs text-black">
									<span className="text-gray-400">z</span> <strong className="font-bold">{z}</strong>
								</div>
							</div>
						</div>
						<div>
							<div className="flex flex-row items-center justify-center px-3 py-3 text-center bg-white m-auto gap-3 rounded-lg shadow-md">
								<div>
									<a href="/canvas/0/1/0/0">
										<img src="/center.svg" alt="location" className="w-4 h-4" />
									</a>
								</div>
							</div>
						</div>
					</div>
				</div>
				<div className="fixed top-0 right-0 bg-white w-fit h-full shadow z-50">
					<Palette
						choices={colorChoices}
						value={canvasState.state.currentBrushColor}
						eraserSelected={canvasState.state.eraserSelected}
						onChange={(newColor) => {
							if (newColor) {
								canvasState.setBrushColor(newColor);
								canvasState.setEraser(false);
							}
						}}
						onEraserClick={() => {
							canvasState.toggleEraser();
						}}
					/>
				</div>
			</ErrorBoundary>
		</div>
	);
};

export default App;
