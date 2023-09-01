import { RGBA, RGBAEquals } from "../types";
import { Color } from "@pixi/core";
import { useHotkeys } from "react-hotkeys-hook";
import { useInstruments } from "../stores/instruments.store";
import { paletteColorsAtom } from "../stores/jotai";
import { useAtomValue } from "jotai";

type PaletteProps = {
	choices: RGBA[];
	value?: RGBA | null;
	eraserSelected: boolean;
	onChange: (color: RGBA | null) => void;
	onEraserClick: () => void;
};

export const Palette = (props: PaletteProps) => {
	const { choices, value, onChange = () => {}, onEraserClick = () => {}, eraserSelected = false } = props;
	const instruments = useInstruments();
	const colorChoices = useAtomValue(paletteColorsAtom);

	// keybindings
	useHotkeys(`e`, () => instruments.setMode("eraser"), []);
	useHotkeys(
		`1`,
		() => {
			instruments.setBrushColor(colorChoices[0]);
			instruments.setMode(`brush`);
		},
		[]
	);
	useHotkeys(
		`2`,
		() => {
			instruments.setBrushColor(colorChoices[1]);
			instruments.setMode(`brush`);
		},
		[]
	);
	useHotkeys(
		`3`,
		() => {
			instruments.setBrushColor(colorChoices[2]);
			instruments.setMode(`brush`);
		},
		[]
	);

	return (
		<ul className=" h-7 flex flex-col flex-nowrap gap-0 text-center">
			{choices.map((choice, i) => {
				const color = new Color(choice);
				return (
					<li key={i}>
						<button
							type="button"
							className="block w-14 h-14"
							style={{ backgroundColor: color.toRgbaString() }}
							onClick={() => onChange(JSON.parse(JSON.stringify(choice)) as RGBA)}
						>
							{!eraserSelected && RGBAEquals(value || null, choice) && `X`}
						</button>
					</li>
				);
			})}
			<li className=" self-end">
				<button
					className="block w-14 h-14"
					style={{ backgroundColor: `#eee` }}
					type="button"
					onClick={onEraserClick}
				>
					{eraserSelected && `X`}
					{!eraserSelected && <img src="/eraser.svg" alt="eraser" className=" block w-full p-3" />}
				</button>
			</li>
		</ul>
	);
};
