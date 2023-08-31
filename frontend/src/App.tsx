import CanvasSingle from "./routes/canvas.single";
import ErrorPage from "./routes/error-page";
import Root from "./routes/root";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

const router = createBrowserRouter([
	{
		path: "/",
		element: <Root />,
		errorElement: <ErrorPage />,
	},
	{
		path: "/canvas/:id",
		element: <CanvasSingle />,
		children: [
			{
				path: "",
				element: <CanvasSingle />,
			},
			{
				path: ":z/:x/:y",
				element: <CanvasSingle />,
			},
		],
	},
]);

export const App = () => {
	return <RouterProvider router={router} />;
};

export default App;
