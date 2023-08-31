import { useEffect, useState } from "react";
import { Point } from "../stores/canvas.store";
import { PixelDrawn } from "../types";

type usePixelPollingProps = {
	canvasId: number;
	topLeft: Point;
	bottomRight: Point;
	from: Date;
	callback: (data: any) => void;
	delayMs: number;
};

const getEndpointFromProps = (props: { canvasId: number; topLeft: Point; bottomRight: Point; from: Date }) => {
	const endpoint = new URL(`http://localhost:1001/poll`);
	endpoint.searchParams.append(`cid`, props.canvasId.toString());
	endpoint.searchParams.append(`from`, props.from.toISOString());
	endpoint.searchParams.append(`tlx`, props.topLeft.x.toString());
	endpoint.searchParams.append(`tly`, props.topLeft.y.toString());
	endpoint.searchParams.append(`brx`, props.bottomRight.x.toString());
	endpoint.searchParams.append(`bry`, props.bottomRight.y.toString());
	return endpoint.toString();
};

export const usePixelPolling = (props: usePixelPollingProps) => {
	const { callback = () => {}, delayMs } = props;

	const [error, setError] = useState<string | null>(null);
	const [isPolling, setIsPolling] = useState<boolean>(false);
	const [endpoint, setEndpoint] = useState<string>(getEndpointFromProps(props));
	const [latestData, setLatestData] = useState<PixelDrawn[]>([]);

	useEffect(() => {
		console.log(`poller props`, props.topLeft.x, props.topLeft.y, props.bottomRight.x, props.bottomRight.y);

		const newEndpoint = getEndpointFromProps({
			canvasId: props.canvasId,
			topLeft: props.topLeft,
			bottomRight: props.bottomRight,
			from: props.from,
		});

		console.log(`newEndpoint`, newEndpoint);

		setEndpoint(newEndpoint);
	}, [props.canvasId, props.topLeft, props.bottomRight, props.from]);

	async function poll() {
		console.log(`Polling ${endpoint}`);

		try {
			const res = await fetch(endpoint);
			if (!res.ok) {
				throw new Error(`HTTP error! status: ${res.status}`);
			}

			const data = await res.json();

			setLatestData(data);
			callback(data);
			setError(null);
		} catch (error) {
			console.error(error);
			setError(error as any);
		}
	}

	function startPolling() {
		setIsPolling(true);
		poll();
		setInterval(poll, delayMs);
	}

	function stopPolling() {
		setIsPolling(false);
	}

	return {
		latestData,
		isPolling,
		error,
		startPolling,
		stopPolling,
	};
};
