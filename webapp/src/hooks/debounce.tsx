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
    }, [searchTermState, searchFunc]);

    return [searchTermState, setSearchTermState] as const;
};
