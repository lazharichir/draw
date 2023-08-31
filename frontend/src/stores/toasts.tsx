import toast from "react-hot-toast";

export const toastFail = (message: string) => {
	toast.error(message, {
		duration: 5000,
	});
};

export const toastSuccess = (message: string) => {
	toast.success(message, {
		duration: 3000,
	});
};

export const toastInfo = (message: string) => {
	toast(message, {
		duration: 1000,
	});
};

export const toastCustom = toast.custom;
