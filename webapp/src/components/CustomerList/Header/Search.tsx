import {TextInput} from '@mantine/core';
import React from 'react';
import {MagnifyIcon} from '@mattermost/compass-icons/components';

import {SetStateDispatch} from '@/types/react';
import {useDebounceSearch} from '@/hooks/debounce';

type Params = {
    setSearchTerm: SetStateDispatch<string>;
}

const SearchBar = ({
    setSearchTerm,
}: Params) => {
    const [searchTermState, setSearchTermState] = useDebounceSearch(setSearchTerm);
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
