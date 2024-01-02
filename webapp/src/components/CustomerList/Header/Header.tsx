import {ActionIcon, Title} from '@mantine/core';
import React, {Dispatch} from 'react';
import styled from 'styled-components';

import {SortAscendingIcon} from '@mattermost/compass-icons/components';

import {CustomerSortOptions, SortDirection} from '@/types/customers';

import {FilterButton} from './FilterButton';

const HeaderContainer = styled.div`
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 7px 1em;
    border-bottom: 1px solid rgba(var(--center-channel-color-rgb), 0.08);
    width: 100%;
`;

type HeaderParams = {
    sortBy: CustomerSortOptions;
    setSortBy: Dispatch<React.SetStateAction<CustomerSortOptions>>
    orderBy: SortDirection;
    setOrderBy: Dispatch<React.SetStateAction<SortDirection>>
}
const Header = ({
    setOrderBy,
    setSortBy,
    sortBy,
    orderBy,
}: HeaderParams) => {
    return (
        <HeaderContainer>
            <Title order={1}>{'Customers'}</Title>
            <div
                style={{
                    display: 'flex',
                    alignItems: 'flex-end',
                    flexDirection: 'row',
                    gap: '.5em',
                }}
            >

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
                <FilterButton
                    setSortBy={setSortBy}
                    sortBy={sortBy}
                />
            </div>

        </HeaderContainer>
    );
};

export {
    Header,
};
