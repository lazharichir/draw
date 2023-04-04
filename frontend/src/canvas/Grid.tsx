import { useCallback } from "react";
import { Graphics } from "@pixi/react";
import * as PIXI from "pixi.js";
import { Viewport } from "pixi-viewport";

interface GridProps {
	color?: [number, number, number];
	lineThickness?: number;
	viewport: Viewport | null;
}

export const Grid = ({ viewport }: GridProps) => {
	if (!viewport) return <></>;

	const draw = (g: PIXI.Graphics) => {
		const side = 1024;
		const offsetX = 0;
		const offsetY = 0;
		const scaled = Math.round(viewport.scaled);
		const x = Math.floor(viewport.center.x) + offsetX;
		const y = Math.floor(viewport.center.y) + offsetY;
		const width = viewport.screenWidth;
		const height = viewport.screenHeight;

		console.log(`Grid.draw`, {
			viewport,
			scaled,
			x,
			y,
			width,
			height,
		});

		g.clear();

		g.lineStyle({
			color: `red`,
			width: 1,
			texture: PIXI.Texture.WHITE,
		});

		// draw horizontal lines
		for (let i = y; i < y + height; i += scaled) {
			g.moveTo(x, i);
			g.lineTo(x + width, i);
		}
	};

	return <Graphics anchor={0.5} draw={draw} />;
};
