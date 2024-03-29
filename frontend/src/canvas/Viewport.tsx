import React, { forwardRef } from "react";
import { Point } from "@pixi/core";
import { Application } from "@pixi/app";
import { PixiComponent, useApp } from "@pixi/react";
import { IClampZoomOptions, Viewport as PixiViewport } from "pixi-viewport";
import { EventSystem, FederatedEvent } from "@pixi/events";

export type ViewportClickedEvent = {
	screen: Point;
	world: Point;
	viewport: PixiViewport;
};

export type ViewportZoomedEndEvent = {
	viewport: PixiViewport;
};

export type ViewportMovedEvent = {
	original?: Point;
	viewport: PixiViewport;
	type:
		| "wheel"
		| "pinch"
		| "animate"
		| "ensureVisible"
		| "snap"
		| "mouse-edges"
		| "follow"
		| "drag"
		| "decelerate"
		| "clamp-x"
		| "clamp-y"
		| "bounce-x"
		| "bounce-y";
};

export type ViewportInitedEvent = {
	viewport: PixiViewport;
};

export type ViewportDragEvent = {
	event: FederatedEvent;
	viewport: PixiViewport;
	screen: Point;
	world: Point;
};

export interface ViewportProps {
	screenWidth: number;
	screenHeight: number;
	worldWidth: number;
	worldHeight: number;
	children?: React.ReactNode;
	// event handlers
	onZoomedEnd?: (e: ViewportZoomedEndEvent) => void;
	onClicked?: (e: ViewportClickedEvent) => void;
	onMoved?: (e: ViewportMovedEvent) => void;
	onInited?: (e: ViewportInitedEvent) => void;
	onPointerMoved?: (e: Point) => void;
	onDragStart?: (e: ViewportDragEvent) => void;
	onDragEnd?: (e: ViewportDragEvent) => void;
	// plugin options
	clampZoomOptions?: IClampZoomOptions;
}

export interface PixiComponentViewportProps extends ViewportProps {
	app: Application;
	clampZoomOptions?: IClampZoomOptions;
}

const PixiComponentViewport = PixiComponent("Viewport", {
	create: (props: PixiComponentViewportProps) => {
		const { screenWidth, screenHeight, worldWidth, worldHeight, app } = props;
		const { ticker, renderer } = app;
		const { events } = renderer;

		const clampZoomOptions = props.clampZoomOptions || {};

		events.cursorStyles.square = "url(/brush.png) 0 0, auto";
		events.cursorStyles.crosshair = "crosshair";

		events.domElement = props.app.view as HTMLCanvasElement;

		const viewport = new PixiViewport({
			screenHeight,
			screenWidth,
			worldWidth,
			worldHeight,
			// ticker,
			events,
			threshold: 5,
			// disableOnContextMenu: true,
		});

		viewport.setZoom(1, true);

		viewport.on(`clicked`, (e) => {
			if (props.onClicked) {
				props.onClicked({
					screen: e.screen,
					world: e.world,
					viewport: e.viewport,
				});
			}
		});

		viewport.on(`zoomed-end`, (e) => {
			if (props.onZoomedEnd) {
				props.onZoomedEnd({
					viewport: viewport,
				});
			}
		});

		viewport.on(`drag-start`, (e) => {
			if (props.onDragStart) {
				props.onDragStart({
					...e,
				});
			}
		});

		viewport.on(`drag-end`, (e) => {
			if (props.onDragEnd) {
				props.onDragEnd({
					...e,
				});
			}
		});

		viewport.on(`moved`, (e) => {
			if (props.onMoved) {
				props.onMoved({
					original: e.original,
					viewport: viewport,
					type: e.type,
				});
			}
		});

		viewport.on(`pointermove`, (e) => {
			if (props.onPointerMoved) {
				props.onPointerMoved(viewport.toWorld(e.global.x, e.global.y));
			}
		});

		viewport
			.drag({})
			.decelerate()
			.pinch({
				percent: 1,
			})
			.wheel({
				center: null,
				percent: 1,
			})
			.clampZoom(clampZoomOptions);

		viewport.resize(window.innerWidth, window.innerHeight);
		viewport.moveCenter(0, 0);

		if (props.onInited) {
			props.onInited({ viewport });
		}

		return viewport;
	},
});

// create a component that can be consumed
// that automatically pass down the app
export const Viewport = forwardRef<PixiViewport, ViewportProps>((props, ref) => {
	const app = useApp();

	// Install EventSystem, if not already (PixiJS 6 doesn't add it by default)
	if (!("events" in app.renderer)) {
		// @ts-ignore
		app.renderer.addSystem(EventSystem, "events");
	}

	return <PixiComponentViewport ref={ref} app={app} {...props} />;
});
