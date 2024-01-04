import {TextInput, createStyles} from '@mantine/core';
import React from 'react';

import {useDebounceSearch} from '@/hooks/debounce';

const useStyles = createStyles(() => ({
    root: {
        width: '100%',
    },
    label: {
        fontSize: '14px',
        fontWeight: 600,
    },
    input: {
        height: '34px',
        color: 'var(--center-channel-color)',
        background: 'var(--center-channel-bg)',
        border: '1px solid #ccc',
        padding: '6px 12px',
        borderColor: 'rgba(var(--center-channel-color-rgb), 0.16)',
        borderRadius: '4px',
        lineHeight: '1.42857143',
        transition: 'border-color ease-in-out .15s, box-shadow ease-in-out .15s, -webkit-box-shadow ease-in-out .15s',
        '&:focus': {
            borderColor: 'rgba(var(--button-bg-rgb), 1)',
            boxShadow: '0 0 0 1px var(--button-bg)',
        },
    },
}));

type Params = {
    handleSearchChange: (text: string) => void;
    disabled: boolean
}
const ConfigSearchBar = ({
    handleSearchChange,
    disabled,
}: Params) => {
    const {classes} = useStyles();

    const [, setSearchTermState] = useDebounceSearch(handleSearchChange);

    return (
        <TextInput
            bg={'transparent'}
            classNames={classes}
            placeholder={'Search for config values'}
            disabled={disabled}
            onChange={(e) => setSearchTermState(e.target.value)}
        />
    );
};

export {
    ConfigSearchBar,
};
