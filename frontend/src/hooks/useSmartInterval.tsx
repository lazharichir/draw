import { useState } from "react";

export type SmartIntervalProps = {
	min: number;
	max: number;
	errorMultiplier: number;
	quietMultiplier: number;
	successMultiplier: number;
	quietMax: number;
	initialValue: number;
};

export const useSmartInterval = (props: SmartIntervalProps) => {
	const [value, setValue] = useState(props.initialValue);

	function success() {
		const nextValue = Math.floor(value * props.successMultiplier);
		if (nextValue < props.min) {
			setValue(props.min);
		} else if (nextValue > props.max) {
			setValue(props.max);
		} else {
			setValue(nextValue);
		}
	}

	function error() {
		if (value === props.max) {
			return;
		}

		const nextValue = Math.ceil(value * props.errorMultiplier);
		if (nextValue > props.max) {
			setValue(props.max);
		} else {
			setValue(nextValue);
		}
	}

	function quiet() {
		if (value === props.max) {
			return;
		}

		const nextValue = Math.floor(value * props.quietMultiplier);

		if (nextValue > props.quietMax) {
			setValue(props.quietMax);
			return;
		} else if (nextValue === props.quietMax) {
			return;
		} else {
			setValue(nextValue);
		}
	}

	return {
		value,
		success,
		error,
		quiet,
	};
};
