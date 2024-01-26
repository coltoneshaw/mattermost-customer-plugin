export function parseJson<T, >(jsonString: string): T | undefined {
    try {
        const parsedData = JSON.parse(jsonString);
        return parsedData;
    } catch (error) {
        // eslint-disable-next-line no-undefined
        return undefined;
    }
}
