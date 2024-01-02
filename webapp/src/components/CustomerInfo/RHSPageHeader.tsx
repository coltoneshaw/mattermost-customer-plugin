import {Title} from '@mantine/core';
import React from 'react';
import styled from 'styled-components';
import {useLocation} from 'react-router-dom';

import {MenuButton} from './MenuButton';

const HeaderContainer = styled.div`
    display: flex;
    align-items: center;
    justify-content: flex-start;
    padding: 7px 1em;
    border-bottom: 1px solid rgba(var(--center-channel-color-rgb), 0.08);
    width: 100%;
    gap: 8px;
    height: 50px;
`;

type Params = {
    id: string;
}
const RhsPageHeader = ({
    id,
}: Params) => {
    const location = useLocation();
    const {pathname} = location;

    let title = pathname.split('/').pop() || 'Customer Info';

    switch (title) {
    case 'packet':
        title = 'Support Packet';
        break;
    case 'config':
        title = 'Config';
        break;
    case 'plugins':
        title = 'Plugins';
        break;
    default:
        title = 'Customer Info';
        break;
    }

    return (
        <HeaderContainer>
            <Title
                sx={{
                    color: 'var(--center-channel-color)',
                    maxWidth: '80%',
                    fontSize: '1em',
                    whiteSpace: 'nowrap',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                }}
            >{title}</Title>
            <MenuButton id={id}/>
        </HeaderContainer>
    );
};

export {
    RhsPageHeader,
};
