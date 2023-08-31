import { InstrumentsProvider } from "./stores/instruments.store";
import { Toaster } from "react-hot-toast";

export const Providers = ({ children }: React.PropsWithChildren<{}>) => {
	return (
		<InstrumentsProvider>
			<Toaster />
			{children}
		</InstrumentsProvider>
	);
};
