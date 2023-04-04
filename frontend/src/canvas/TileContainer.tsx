import { Container, Sprite } from "@pixi/react";
import { TileData } from "../types";

export const TileContainer = (props: { side: number; tiles: TileData[] }) => {
	const { side, tiles } = props;
	return (
		<Container>
			{tiles.map((tile: TileData) => {
				return (
					<Sprite
						key={`${tile.x}x${tile.y}_${tile.anchor}`}
						image={`http:\/\/localhost:1001/tile/${tile.x}x${tile.y}_${side}.png`}
						{...tile}
					/>
				);
			})}
		</Container>
	);
};
