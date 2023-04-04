import { TileData } from "../types";
import uniqWith from "lodash.uniqwith";

// An enum with all the types of actions to use in our reducer
enum ActionKind {
	ADD = "ADD",
}

// An interface for our actions
interface Action {
	type: `ADD`;
	payload: TileData[];
}

// An interface for our state
interface State {
	tiles: TileData[];
}

// Our reducer function that uses a switch statement to handle our actions
export function useTilesReducer(state: State, action: Action): State {
	const { type, payload } = action;
	switch (type) {
		case `ADD`:
			if (payload.length === 0) return state;

			const uniqueTiles = uniqWith(
				[...state.tiles, ...payload] as TileData[],
				(a, b) => a.x === b.x && a.y === b.y
			);

			// console.log(
			// 	`> uniqueTiles`,
			// 	uniqueTiles.map((t) => `${t.x}x${t.y}`)
			// );

			return {
				tiles: uniqueTiles,
			};
		default:
			return state;
	}
	return state;
}
