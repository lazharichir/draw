import { Color } from "@pixi/core";
import { RGBA } from "./types";

export const putPixelRemotely = async (canvasID: number, x: number, y: number, pixelColor: RGBA) => {
    try {
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
            console.error(`Error putting pixel remotely`, res);
        }
    } catch (error) {
        console.error(`Error putting pixel remotely [catch]`, error);
    }
};

export const erasePixelRemotely = async (canvasID: number, x: number, y: number) => {
    try {
        const res = await fetch(`http://localhost:1001/pixel/${canvasID}/${x}/${y}`, {
            method: `DELETE`,
            headers: {
                "Content-Type": `application/json`,
            },
        });

        console.log(`erasePixelRemotely.res`, res.ok, res);

        if (!res.ok) {
            console.error(`Error erasing pixel remotely`, res);
        } else {
            console.log(`pixel erased remotely`);
        }
    } catch (error) {
        console.error(`Error erasing pixel remotely [catch]`, error);
    }
};