import {TextInput} from '@mantine/core';
import React, {useCallback, useEffect, useState} from 'react';
import {MagnifyIcon} from '@mattermost/compass-icons/components';
import debounce from 'debounce';

import {SetStateDispatch} from '@/types/react';

type Params = {
    setSearchTerm: SetStateDispatch<string>;
}

const SearchBar = ({
    setSearchTerm,
}: Params) => {
    const [searchTermState, setSearchTermState] = useState<string>('');

    // eslint-disable-next-line react-hooks/exhaustive-deps
    const debouncedSearchTerm = useCallback(
        debounce((searchTerm: string) => setSearchTerm(searchTerm), 300),
        [setSearchTerm],
    );

    useEffect(() => {
        debouncedSearchTerm(searchTermState);

        return () => {
            debouncedSearchTerm.clear();
        };
    }, [searchTermState, debouncedSearchTerm]);

    return (
        <TextInput
            bg={'transparent'}
            size='lg'
            sx={{
                width: '100%',
                input: {
                    color: 'rgba(var(--center-channel-color-rgs))',
                    background: 'transparent',
                },
            }}
            value={searchTermState}
            onChange={(event) => setSearchTermState(event.currentTarget.value)}
            placeholder='Search'
            icon={<MagnifyIcon size='18px'/>}
        />
    );
};

export {
    SearchBar,
};
