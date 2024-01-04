import debounce from 'debounce';
import {useEffect, useState} from 'react';

export const useDebounceSearch = (searchFunc: (term: string) => void) => {
    const [searchTermState, setSearchTermState] = useState('');

    useEffect(() => {
        const debouncedSearch = debounce(() => searchFunc(searchTermState), 300);
        debouncedSearch();

        return () => {
            debouncedSearch.clear();
        };

    // adding searchFunc to the dependency array causes an infinite loop
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [searchTermState]);

    return [searchTermState, setSearchTermState] as const;
};
