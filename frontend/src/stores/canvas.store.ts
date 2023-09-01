import { useReducer } from 'react';
import { PixelSnapshot, RGBA, TileData } from "../types"
import uniqWith from 'lodash.uniqwith';
import { Viewport } from 'pixi-viewport';
import { retile } from '../canvas/retile';

export type Point = {
    x: number
    y: number
}

export type CanvasState = {
    canvasId: number
    screenWidth: number
    screenHeight: number
    worldWidth: number
    worldHeight: number
    backgroundColor: RGBA
    gridColor: RGBA
    side: number
    center: Point
    lastClick?: Point
    scale: number
    tiles: TileData[],
    pixels: PixelSnapshot,
}

export type SetPixelData = {
    id: string
    x: number
    y: number
    color: RGBA
    at: number
    deleted: boolean
}

export type CanvasAction =
    | { type: 'SET_CANVAS_ID', id: number }
    | { type: 'SET_SCREEN_SIZE', width: number, height: number }
    | { type: 'SET_WORLD_SIZE', width: number, height: number }
    | { type: 'SET_BACKGROUND_COLOR', color: RGBA }
    | { type: 'SET_GRID_COLOR', color: RGBA }
    | { type: 'SET_SIDE', side: number }
    | { type: 'SET_CENTER', x: number, y: number }
    | { type: 'SET_SCALE', scale: number }
    | { type: 'ADD_TILES', tiles: TileData[] }
    | { type: 'SET_PIXELS', pixels: SetPixelData[] }
    | { type: 'UNSET_PIXELS', pixels: SetPixelData[] }
    | { type: 'SET_LAST_CLICK', x: number, y: number }
    | { type: 'RETILE', viewport: Viewport }

export const canvasReducer = (state: CanvasState, action: CanvasAction): CanvasState => {
    switch (action.type) {
        case 'RETILE':
            return {
                ...state,
                tiles: uniqWith(
                    [
                        ...state.tiles,
                        ...retile(action.viewport, state.tiles, state.side),
                    ],
                    (a, b) => a.x === b.x && a.y === b.y
                ),
            }
        case 'SET_CANVAS_ID':
            return {
                ...state,
                canvasId: action.id
            }
        case 'SET_SCREEN_SIZE':
            return {
                ...state,
                screenWidth: action.width,
                screenHeight: action.height
            }
        case 'SET_WORLD_SIZE':
            return {
                ...state,
                worldWidth: action.width,
                worldHeight: action.height
            }
        case 'SET_BACKGROUND_COLOR':
            return {
                ...state,
                backgroundColor: action.color
            }
        case 'SET_GRID_COLOR':
            return {
                ...state,
                gridColor: action.color
            }
        case 'SET_SIDE':
            return {
                ...state,
                side: action.side
            }
        case 'SET_CENTER':
            return {
                ...state,
                center: {
                    x: action.x,
                    y: action.y
                }
            }
        case 'SET_LAST_CLICK':
            return {
                ...state,
                lastClick: {
                    x: action.x,
                    y: action.y
                }
            }
        case 'SET_SCALE':
            return {
                ...state,
                scale: action.scale
            }
        case 'ADD_TILES':
            return {
                ...state,
                tiles: uniqWith([...state.tiles, ...action.tiles], (a, b) => a.x === b.x && a.y === b.y)
            }
        case 'SET_PIXELS':
            {
                let pixels = { ...state.pixels }
                for (const data of action.pixels) {
                    const { x, y, at, color, deleted, id } = data

                    pixels[x] = pixels[x] || {}
                    pixels[x][y] = pixels[x][y] || []

                    const newPixel = {
                        id: id,
                        at: at,
                        color: color,
                        erased: deleted,
                    }

                    // console.log(`set pixel at: `, x, y)

                    pixels[x][y] = [
                        newPixel,
                        ...pixels[x][y],
                    ]
                }

                // console.log(`nextState.pixels: `, pixels)
                return { ...state, pixels: pixels }
            }
        case 'UNSET_PIXELS':
            {
                let pixels = { ...state.pixels }
                for (const data of action.pixels) {
                    pixels = ensureXYPath(pixels, data.x, data.y)
                    pixels[data.x][data.y] = [...pixels[data.x][data.y].filter(item => item.id !== data.id)]
                }
                return { ...state, pixels }
            }
        default:
            return state;
    }
}

function ensureXYPath(pixels: PixelSnapshot, x: number, y: number): PixelSnapshot {
    if (pixels[x] === undefined) {
        pixels[x] = {}
    }
    if (pixels[x][y] === undefined) {
        pixels[x][y] = []
    }
    return pixels
}

export const useCanvasReducer = (initialValue: Partial<CanvasState>) => {
    const [state, dispatch] = useReducer(canvasReducer, { ...initialState, ...initialValue });
    return {
        state,
        setCanvasId: (id: number) => dispatch({ type: 'SET_CANVAS_ID', id }),
        setScreenSize: (width: number, height: number) => dispatch({ type: 'SET_SCREEN_SIZE', width, height }),
        setWorldSize: (width: number, height: number) => dispatch({ type: 'SET_WORLD_SIZE', width, height }),
        setBackgroundColor: (color: RGBA) => dispatch({ type: 'SET_BACKGROUND_COLOR', color }),
        setGridColor: (color: RGBA) => dispatch({ type: 'SET_GRID_COLOR', color }),
        setSide: (side: number) => dispatch({ type: 'SET_SIDE', side }),
        setCenter: (x: number, y: number) => dispatch({ type: 'SET_CENTER', x, y }),
        setScale: (scale: number) => dispatch({ type: 'SET_SCALE', scale }),
        addTiles: (tiles: TileData[]) => dispatch({ type: 'ADD_TILES', tiles }),
        setPixels: (pixels: SetPixelData[]) => dispatch({ type: 'SET_PIXELS', pixels }),
        unsetPixels: (pixels: SetPixelData[]) => dispatch({ type: 'UNSET_PIXELS', pixels }),
        setLastClick: (x: number, y: number) => dispatch({ type: 'SET_LAST_CLICK', x, y }),
        retile: (viewport: Viewport) => dispatch({ type: 'RETILE', viewport }),
    };
}

const initialState: CanvasState = {
    canvasId: 0,
    screenWidth: 0,
    screenHeight: 0,
    worldWidth: 9999999999,
    worldHeight: 9999999999,
    backgroundColor: {
        r: 255,
        g: 240,
        b: 229,
        a: 1,
    },
    gridColor: { // d3d3d3 in rgba format: 
        r: 211,
        g: 211,
        b: 211,
        a: 1,
    },
    side: 1024,
    scale: 1,
    center: { x: 0, y: 0 },
    lastClick: undefined,
    tiles: [],
    pixels: {} as PixelSnapshot // Assuming this is an appropriate default value. Adjust as needed.
};
