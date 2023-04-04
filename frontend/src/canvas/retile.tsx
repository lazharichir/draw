import { Viewport } from "pixi-viewport";
import { TileData } from "../types";
import throttle from "lodash.throttle";
import debounce from "lodash.debounce";

export const retile = (vp: Viewport, tiles: TileData[], side: number): TileData[] => {
	const centerX = Math.floor(vp.center.x);
	const centerY = Math.floor(vp.center.y);
	const screenWidthInWorldPixels = Math.ceil(vp.screenWidthInWorldPixels);
	const screenHeightInWorldPixels = Math.ceil(vp.screenHeightInWorldPixels);

	// split the viewport into SIDE-sized tiles with modulo = 0 as cut points
	// and round them to the nearest SIDE (e.g. 1024)
	// this way we can easily find the tiles that are in the viewport and also the ones around it to add them

	const xStartRounded = Math.floor(centerX / side) * side - side;
	const xEndRounded = Math.ceil((centerX + screenWidthInWorldPixels) / side) * side + side;

	const yStartRounded = Math.floor(centerY / side) * side - side;
	const yEndRounded = Math.ceil((centerY + screenHeightInWorldPixels) / side) * side + side;

	// init new tiles array
	const newTiles: TileData[] = [];

	// add all sprites that are in the viewport
	for (let x = xStartRounded; x <= xEndRounded; x += side) {
		for (let y = yStartRounded; y <= yEndRounded; y += side) {
			const sprite = tiles.find((s) => s.x === x && s.y === y);

			// console.log(
			// 	`CURRENT [${x}x${y}] ${sprite ? `[X]` : `[ ]`} [${tiles.map(
			// 		(s) => `${s.x}x${s.y}`
			// 	)}]`
			// );

			if (!sprite) {
				// console.log(`> adding sprite`, x, y);
				newTiles.push({
					x,
					y,
					anchor: 0,
				});
			}
		}
	}

	return newTiles;
};

export const throttledRetile = throttle(retile, 200);

export const debouncedRetile = debounce(retile, 200);
