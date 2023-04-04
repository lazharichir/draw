import { useState, useEffect } from "react";

export const useResize = () => {
	const [screenSize, setScreenSize] = useState([
		window.innerWidth,
		window.innerHeight,
	]);

	useEffect(() => {
		const onResize = () => {
			requestAnimationFrame(() => {
				console.log(
					"resized to: ",
					window.innerWidth,
					"x",
					window.innerHeight
				);
				setScreenSize([window.innerWidth, window.innerHeight]);
			});
		};

		window.addEventListener("resize", onResize);

		return () => {
			window.removeEventListener("resize", onResize);
		};
	}, []);

	return screenSize;
};
