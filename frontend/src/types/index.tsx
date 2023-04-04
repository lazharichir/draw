export type TileData = {
	x: number;
	y: number;
	anchor: number;
};

export const RGBAEquals = (a: RGBA | null, b: RGBA | null) => {
	if (a !== null && b === null) return false;
	if (a === null && b !== null) return false;
	if (a === null && b === null) return true;
	a = a as RGBA;
	b = b as RGBA;
	return a.r === b.r && a.g === b.g && a.b === b.b && a.a === b.a;
};

export const RGBADebugStr = (a: RGBA | null | undefined) => {
	if (!a) return `rgba(null)`;
	return `rgba(${a.r}, ${a.g}, ${a.b}, ${a.a})`;
};

export type RGBA = {
	r: number;
	g: number;
	b: number;
	a: number;
};

export type PixelDrawn = {
	at: number;
	x: number;
	y: number;
	color: RGBA;
};

export type ErasedPixel = {
	at: number;
	x: number;
	y: number;
};

export type PixelSnapshot = Record<number, Record<number, PixelSnapshotItem>>;

export type PixelSnapshotItem = {
	at: number;
	color: RGBA;
	erased: boolean;
};
