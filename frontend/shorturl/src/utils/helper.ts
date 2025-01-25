export function debounce(callback: (...args: never[]) => unknown, delay: number) {
    let timeout: number | undefined;

    return (...args: never[]) => {
        clearTimeout(timeout)
        timeout = setTimeout(() => {
            callback(...args)
        }, delay)
    };
}