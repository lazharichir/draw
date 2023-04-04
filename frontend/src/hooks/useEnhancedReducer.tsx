import {
	Reducer,
	ReducerAction,
	ReducerState,
	useCallback,
	useMemo,
	useReducer,
	useRef,
} from "react";

export const useEnhancedReducer = <R extends Reducer<any, any>>(
	reducer: R,
	initState: ReducerState<R>,
	initializer?: Parameters<typeof useReducer>[2],
	middlewares: Array<
		(
			state: ReducerState<R>
		) => (
			getState: () => ReducerState<R>
		) => (
			next: (action: ReducerAction<R>) => any
		) => (action: ReducerAction<R>) => any
	> = []
) => {
	const lastState = useRef(initState);
	const getState = useCallback(() => lastState.current, []);
	const enhancedReducer = useRef(
		(state: ReducerState<R>, action: ReducerAction<R>) =>
			(lastState.current = reducer(state, action))
	).current; // to prevent reducer called twice, per: https://github.com/facebook/react/issues/16295
	const [state, dispatch] = useReducer(
		enhancedReducer,
		initState,
		initializer
	);
	const middlewaresRef = useRef(middlewares);
	//use useMemo instead of useRef to avoid redundant calculation
	const enhancedDispatch = useMemo(
		() =>
			middlewaresRef.current.reduceRight(
				(acc, mdw) => (action) => mdw(state)(getState)(acc)(action),
				dispatch
			),
		[]
	);
	return [state, enhancedDispatch, getState];
};
