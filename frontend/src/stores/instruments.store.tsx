import { MutableRefObject, PropsWithChildren, createContext, useContext, useReducer, useRef } from "react";
import { RGBA } from "../types";

type InstrumentContextType = {
	state: InstrumentState;
	dispatch: React.Dispatch<Action>;
	setMode: (mode: `brush` | `eraser`) => void;
	setBrushSize: (size: number) => void;
	setBrushColor: (color: RGBA) => void;
	current: () => InstrumentState;
};

export type InstrumentState = {
	mode: `brush` | `eraser`;
	brushSize: number;
	brushColor: RGBA;
};

const initialState = (): InstrumentState => ({
	mode: `brush`,
	brushSize: 1,
	brushColor: { r: 0, g: 0, b: 0, a: 1 },
});

export const InstrumentContext = createContext<InstrumentContextType>({
	state: initialState(),
	current: () => initialState(),
	dispatch: () => {},
	setMode: () => {},
	setBrushSize: () => {},
	setBrushColor: () => {},
});

type Action =
	| { type: `SET_MODE`; mode: `brush` | `eraser` }
	| { type: `SET_BRUSH_SIZE`; size: number }
	| { type: `SET_BRUSH_COLOR`; color: RGBA };

const reducer = (state: InstrumentState, action: Action): InstrumentState => {
	switch (action.type) {
		case `SET_MODE`:
			return { ...state, mode: action.mode };
		case `SET_BRUSH_SIZE`:
			return { ...state, brushSize: action.size };
		case `SET_BRUSH_COLOR`:
			return { ...state, brushColor: action.color };
	}
};

export const InstrumentsProvider: React.FC<PropsWithChildren> = ({ children }) => {
	const [state, dispatch] = useReducer(reducer, initialState());
	const ref = useRef<InstrumentState>(state);
	ref.current = state;

	const setMode = (mode: `brush` | `eraser`) => dispatch({ type: `SET_MODE`, mode });

	const setBrushSize = (size: number) => dispatch({ type: `SET_BRUSH_SIZE`, size });

	const setBrushColor = (color: RGBA) => dispatch({ type: `SET_BRUSH_COLOR`, color });

	function current() {
		return ref.current;
	}

	const value = { state, dispatch, current, setMode, setBrushSize, setBrushColor };

	return <InstrumentContext.Provider value={value}>{children}</InstrumentContext.Provider>;
};

export const useInstruments = () => useContext(InstrumentContext);
