import * as PIXI from "pixi.js";
import { useMemo } from "react";
import { PixelSnapshot, RGBA } from "../types";
import { Sprite, Container } from "@pixi/react";
import { OutlineFilter } from "@pixi/filter-outline";

export const ClientPixelContainer = (props: { snapshot: PixelSnapshot; backgroundColor: RGBA }) => {
	const { snapshot, backgroundColor } = props;

	const pixels: JSX.Element[] = useMemo(() => {
		const els: JSX.Element[] = [];
		let k = 0;

		Object.keys(snapshot).forEach((x) => {
			Object.keys(snapshot[+x]).forEach((y) => {
				const { color, erased } = snapshot[+x][+y];
				els.push(
					<Sprite
						key={k++}
						texture={PIXI.Texture.WHITE}
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
	}, [snapshot, backgroundColor]);

	return <Container>{pixels}</Container>;
};
