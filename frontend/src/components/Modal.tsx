import React, { useEffect } from "react";
import { Point } from "../stores/canvas.store";

export default function Modal({
	onPointChange = () => {},
	onClose = () => {},
	isOpen = false,
}: {
	isOpen: boolean;
	onPointChange: (point: Point) => void;
	onClose: () => void;
}) {
	const [point, setPoint] = React.useState<Point>({
		x: 0,
		y: 0,
	});

	return (
		<div className={isOpen ? `` : `hidden`}>
			<div className=" justify-center items-center flex overflow-x-hidden overflow-y-auto fixed inset-0 z-50 outline-none focus:outline-none">
				<div className="relative w-auto my-6 mx-auto max-w-3xl">
					{/*content*/}
					<div className="border-0 rounded-lg shadow-lg relative flex flex-col w-full bg-white outline-none focus:outline-none">
						{/*header*/}
						<div className="flex items-start justify-between p-5 border-b border-solid border-slate-200 rounded-t">
							<h3 className="text-3xl font-semibold">Go To Coordinates</h3>
							<button
								className="p-1 ml-auto bg-transparent border-0 text-black opacity-5 float-right text-3xl leading-none font-semibold outline-none focus:outline-none"
								onClick={() => onClose()}
							>
								<span className="bg-transparent text-black opacity-5 h-6 w-6 text-2xl block outline-none focus:outline-none">
									×
								</span>
							</button>
						</div>
						{/*body*/}
						<div className="relative p-6 flex-auto">
							<div>
								<input
									type="number"
									placeholder="x"
									step={1}
									onChange={(e) => setPoint({ ...point, x: parseInt(e.currentTarget.value) })}
								/>
								<input
									type="number"
									placeholder="y"
									step={1}
									onChange={(e) => setPoint({ ...point, y: parseInt(e.currentTarget.value) })}
								/>
							</div>
						</div>
						{/*footer*/}
						<div className="flex items-center justify-end p-6 border-t border-solid border-slate-200 rounded-b">
							<button
								className="text-red-500 background-transparent font-bold uppercase px-6 py-2 text-sm outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
								type="button"
								onClick={() => onClose()}
							>
								Cancel
							</button>
							<button
								className="bg-emerald-500 text-white active:bg-emerald-600 font-bold uppercase text-sm px-6 py-3 rounded shadow hover:shadow-lg outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
								type="button"
								onClick={() => {
									onPointChange(point);
									onClose();
								}}
							>
								Go
							</button>
						</div>
					</div>
				</div>
			</div>
			<div className="opacity-25 fixed inset-0 z-40 bg-black pointer-events-none"></div>
		</div>
	);
}
