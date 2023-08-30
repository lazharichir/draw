import { RGBA, RGBAEquals } from "../types";
import { Color } from "@pixi/core";

type PaletteProps = {
	choices: RGBA[];
	value?: RGBA | null;
	eraserSelected: boolean;
	onChange: (color: RGBA | null) => void;
	onEraserClick: () => void;
};

export const Palette = (props: PaletteProps) => {
	const { choices, value, onChange = () => {}, onEraserClick = () => {}, eraserSelected = false } = props;
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
							onClick={() => onChange(choice)}
						>
							{RGBAEquals(value || null, choice) && `X`}
						</button>
					</li>
				);
			})}
			<li key={`null`}>
				<button
					className="block w-14 h-14"
					style={{ backgroundColor: `#ccc` }}
					type="button"
					onClick={onEraserClick}
				>
					{eraserSelected && `X`}
				</button>
			</li>
		</ul>
	);
};
