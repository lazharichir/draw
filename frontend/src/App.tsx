import * as PIXI from "pixi.js";
import { Canvas } from "./canvas/Canvas";
import { Palette } from "./canvas/Palette";
import { ErrorBoundary } from "./ErrorBoundary";
import { useResize } from "./hooks/useResize";
import { useAtom, useAtomValue, useSetAtom } from "jotai";
import { eraserSelectedAtom, paletteColorsAtom, selectedColorAtom } from "./stores/jotai";

// set some global pixi settings
PIXI.settings.RESOLUTION = window.devicePixelRatio;
PIXI.settings.SCALE_MODE = PIXI.SCALE_MODES.NEAREST;

export const App = () => {
	const [screenWidth, screenHeight] = useResize();
	const colorChoices = useAtomValue(paletteColorsAtom);
	const [eraserSelected, setEraserSelected] = useAtom(eraserSelectedAtom);
	const [selectedColor, setSelectedColor] = useAtom(selectedColorAtom);

	return (
		<div className="fixed top-0 left-0 w-full h-full">
			<ErrorBoundary>
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
				/>
			</ErrorBoundary>
		</div>
	);
};

export default App;
