import { Color } from "@pixi/core";
import { PixelDrawn, RGBA } from "./types";
import { Point } from "./stores/canvas.store";

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

export const pollPixelsRemotely = async (canvasId: number, topLeft: Point, bottomRight: Point, since: Date): Promise<PixelDrawn[]> => {
    // build endpoint
    const endpoint = new URL(`http://localhost:1001/poll`);
    endpoint.searchParams.append(`cid`, canvasId.toString());
    endpoint.searchParams.append(`from`, since.toISOString());
    endpoint.searchParams.append(`tlx`, topLeft.x.toString());
    endpoint.searchParams.append(`tly`, topLeft.y.toString());
    endpoint.searchParams.append(`brx`, bottomRight.x.toString());
    endpoint.searchParams.append(`bry`, bottomRight.y.toString());

    const res = await fetch(endpoint.toString());
    if (!res.ok) {
        throw Error(`Error putting pixel remotely: ${res.statusText}`);
    }

    return await res.json() || [];
};

export const putPixelRemotely = async (canvasID: number, x: number, y: number, pixelColor: RGBA) => {
    const colr = new Color(pixelColor);
    const [r, g, b] = colr.toUint8RgbArray();
    const a = Math.round(colr.alpha * 255);

    const res = await fetch(`http://localhost:1001/pixel/${canvasID}/${x}/${y}/${r}/${g}/${b}/${a}`, {
        method: `PUT`,
        headers: {
            "Content-Type": `application/json`,
        },
    });

    console.log(`putPixelRemotely.res`, res.ok, res);

    if (!res.ok) {
        throw Error(`Error putting pixel remotely: ${res.statusText}`);
    }
};

export const erasePixelRemotely = async (canvasID: number, x: number, y: number) => {
    const res = await fetch(`http://localhost:1001/pixel/${canvasID}/${x}/${y}`, {
        method: `DELETE`,
        headers: {
            "Content-Type": `application/json`,
        },
    });

    console.log(`erasePixelRemotely.res`, res.ok, res);

    if (!res.ok) {
        throw Error(`Error erasing pixel remotely: ${res.statusText}`);
    }
};