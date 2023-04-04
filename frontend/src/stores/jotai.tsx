import { atom } from "jotai";
import { atomWithImmer } from "jotai-immer";
import { RGBA } from "../types";

export const paletteColorsAtom = atomWithImmer<RGBA[]>([
	{ r: 255, g: 0, b: 0, a: 1 },
	{ r: 0, g: 255, b: 0, a: 1 },
	{ r: 0, g: 0, b: 255, a: 1 },
	{ r: 255, g: 255, b: 0, a: 1 },
	{ r: 255, g: 0, b: 255, a: 1 },
	{ r: 0, g: 255, b: 255, a: 1 },
	{ r: 200, g: 200, b: 200, a: 1 },
	{ r: 100, g: 100, b: 100, a: 1 },
	{ r: 0, g: 0, b: 0, a: 1 },
]);

export const selectedColorAtom = atom<RGBA | null>({ r: 255, g: 255, b: 0, a: 1 });

export const eraserSelectedAtom = atom(false);
