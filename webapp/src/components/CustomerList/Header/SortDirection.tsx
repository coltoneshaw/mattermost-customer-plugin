import React from 'react';

import {ActionIcon} from '@mantine/core';

import {SortAscendingIcon} from '@mattermost/compass-icons/components';

import {SortDirection} from '@/types/customers';
import {SetStateDispatch} from '@/types/react';

type Params = {
    orderBy: SortDirection;
    setOrderBy: SetStateDispatch<SortDirection>;
}
const SortDirectionButton = ({
    orderBy,
    setOrderBy,
}: Params) => {
    return (
        <ActionIcon
            sx={{
                height: '32px',
                width: '32px',
                color: 'rgba(var(--center-channel-color-rgb), 0.72)',
                transform: orderBy === SortDirection.DirectionDesc ? 'rotate(180deg)' : '',
                '&:hover': {
                    backgroundColor: 'rgba(var(--center-channel-color-rgb), 0.08)',
                    color: 'rgba(var(--center-channel-color-rgb), 0.72)',
                },
            }}
            onClick={() => {
                if (orderBy === SortDirection.DirectionAsc) {
                    setOrderBy(SortDirection.DirectionDesc);
                } else {
                    setOrderBy(SortDirection.DirectionAsc);
                }
            }}
        >
            <SortAscendingIcon
                size={18}
            />
        </ActionIcon>
    );
};

export {
    SortDirectionButton,
};
