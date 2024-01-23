
/**
 *
 * @param property property to sort the array based on
 * @returns function to be used in `.sort()`
 */
function dynamicSort<T>(property: string): T {
    let sortOrder = 1;
    if (property.startsWith('-')) {
        sortOrder = -1;
        // eslint-disable-next-line no-param-reassign
        property = property.substr(1);
    }

    // @ts-expect-error - Too generic to have correct types.
    return (a: Record<string, string>, b: Record<string, string>) => {
    /* next line works with strings and numbers,
     * and you may want to customize it to your needs
     */

        let result = 0;
        if (a[property] < b[property]) {
            result = -1;
        }
        if (a[property] > b[property]) {
            result = 1;
        }

        // const result = (a[property] < b[property]) ? -1 : (a[property] > b[property]) ? 1 : 0;
        return result * sortOrder;
    };
}

export {
    dynamicSort,
};
