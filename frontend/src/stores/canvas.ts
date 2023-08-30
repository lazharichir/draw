import { useReducer } from 'react';
import { PixelSnapshot, PixelSnapshotItem, RGBA, TileData } from "../types"
import uniqWith from 'lodash.uniqwith';

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
    lastClick: Point
    scale: number
    currentBrushColor: RGBA,
    eraserSelected: boolean,
    tiles: TileData[],
    pixels: PixelSnapshot,
}

export type SetPixelData = {
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
    | { type: 'SET_BRUSH_COLOR', color: RGBA }
    | { type: 'SET_ERASER', active: boolean }
    | { type: 'TOGGLE_ERASER' }
    | { type: 'ADD_TILES', tiles: TileData[] }
    | { type: 'SET_PIXELS', pixels: SetPixelData[] }
    | { type: 'SET_LAST_CLICK', x: number, y: number }

export const canvasReducer = (state: CanvasState, action: CanvasAction): CanvasState => {
    switch (action.type) {
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
        case 'SET_BRUSH_COLOR':
            return {
                ...state,
                currentBrushColor: action.color
            }
        case 'SET_ERASER':
            return {
                ...state,
                eraserSelected: action.active
            }
        case 'TOGGLE_ERASER':
            return {
                ...state,
                eraserSelected: !state.eraserSelected
            }
        case 'ADD_TILES':
            return {
                ...state,
                tiles: uniqWith([...state.tiles, ...action.tiles], (a, b) => a.x === b.x && a.y === b.y)
            }
        case 'SET_PIXELS':
            const next = { ...state }
            for (const data of action.pixels) {
                if (next.pixels[data.x] === undefined) {
                    next.pixels[data.x] = {}
                }
                if (next.pixels[data.x][data.y] === undefined) {
                    next.pixels[data.x][data.y] = {} as PixelSnapshotItem
                }
                next.pixels[data.x][data.y] = {
                    at: data.at,
                    color: data.color,
                    erased: data.deleted,
                }
            }
            return next
        default:
            return state;
    }
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
        setBrushColor: (color: RGBA) => dispatch({ type: 'SET_BRUSH_COLOR', color }),
        toggleEraser: () => dispatch({ type: 'TOGGLE_ERASER' }),
        setEraser: (active: boolean) => dispatch({ type: 'SET_ERASER', active }),
        addTiles: (tiles: TileData[]) => dispatch({ type: 'ADD_TILES', tiles }),
        setPixels: (pixels: SetPixelData[]) => dispatch({ type: 'SET_PIXELS', pixels }),
        setLastClick: (x: number, y: number) => dispatch({ type: 'SET_LAST_CLICK', x, y }),
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
    gridColor: {
        r: 223,
        g: 219,
        b: 217,
        a: 1,
    },
    side: 1024,
    scale: 1,
    center: { x: 0, y: 0 },
    lastClick: { x: 0, y: 0 },
    currentBrushColor: { r: 0, g: 0, b: 0, a: 1 },
    eraserSelected: false,
    tiles: [],
    pixels: {} as PixelSnapshot // Assuming this is an appropriate default value. Adjust as needed.
};
