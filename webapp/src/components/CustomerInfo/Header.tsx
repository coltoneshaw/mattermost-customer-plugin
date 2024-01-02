import {Title} from '@mantine/core';
import React from 'react';
import styled from 'styled-components';
import {ArrowBackIosIcon} from '@mattermost/compass-icons/components';
import {useHistory} from 'react-router-dom';

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
    name: string;
    id: string;
}
export const Header = ({
    name,
    id,
}: Params) => {
    const history = useHistory();
    return (
        <HeaderContainer>
            <span
                style={{
                    color: 'var(--center-channel-color)',
                    fontSize: '1em',
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                }}
                onClick={() => {
                    history.push('/customers');
                }}
            >
                <ArrowBackIosIcon size={16}/>
            </span>
            <Title
                sx={{
                    color: 'var(--center-channel-color)',
                    maxWidth: '80%',
                    fontSize: '1em',
                    whiteSpace: 'nowrap',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                }}
            >{name || 'Customer Info'}</Title>
            <MenuButton id={id}/>
        </HeaderContainer>
    );
};
