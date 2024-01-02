import React from 'react';
import styled from 'styled-components';

import {CustomerSortOptions, SortDirection} from '@/types/customers';

import {SetStateDispatch} from '@/types/react';

import {SortByButton} from './SortByButton';
import {SortDirectionButton} from './SortDirection';
import {SearchBar} from './Search';

const HeaderContainer = styled.div`
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 7px 1em;
    border-bottom: 1px solid rgba(var(--center-channel-color-rgb), 0.08);
    width: 100%;
    gap: 8px;
    height: 50px;
`;

type HeaderParams = {
    sortBy: CustomerSortOptions;
    setSortBy: SetStateDispatch<CustomerSortOptions>;
    orderBy: SortDirection;
    setOrderBy: SetStateDispatch<SortDirection>;
    setSearchTerm:SetStateDispatch<string>;
}
const Header = ({
    setOrderBy,
    setSortBy,
    sortBy,
    orderBy,
    setSearchTerm,
}: HeaderParams) => {
    return (
        <HeaderContainer>
            <SearchBar
                setSearchTerm={setSearchTerm}
            />
            <div
                style={{
                    display: 'flex',
                    alignItems: 'flex-end',
                    flexDirection: 'row',
                    gap: '.5em',
                }}
            >
                <SortDirectionButton
                    orderBy={orderBy}
                    setOrderBy={setOrderBy}
                />
                <SortByButton
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
