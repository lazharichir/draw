import { Viewport } from "pixi-viewport";
import { useState, useRef } from "react";
import { useEnhancedReducer } from "../hooks/useEnhancedReducer";
import { RGBA, PixelSnapshot } from "../types";
import { useTilesReducer } from "./useTilesReducer";

const useSurface = (worldWidth: number, worldHeight: number, brushColor: RGBA) => {
	const [tileState, dispatch, getTileState] = useEnhancedReducer(useTilesReducer, { tiles: [] });
	const [[lastPointerX, lastPointerY], setLastPointer] = useState<number[]>([]);
	const [lastClick, setLastClick] = useState<[number, number] | null>(null);
	const [snapshot, setSnapshot] = useState<PixelSnapshot>({});
	const viewportRef = useRef<Viewport>(null);

	const drawPixel = (x: number, y: number, color: RGBA) => {
		setSnapshot((prev) => ({
			...prev,
			[+x]: {
				...(prev[+x] || {}),
				[+y]: {
					at: Date.now(),
					color: color,
					erased: false,
				},
			},
		}));
	};

	const erasePixel = (x: number, y: number) => {
		setSnapshot((prev) => ({
			...prev,
			[+x]: {
				...(prev[+x] || {}),
				[+y]: {
					at: Date.now(),
					color: brushColor,
					erased: true,
				},
			},
		}));
	};

	return {
		drawPixel,
		erasePixel,
	};
};
