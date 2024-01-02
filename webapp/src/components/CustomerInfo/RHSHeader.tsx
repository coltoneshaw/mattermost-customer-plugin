import {Title} from '@mantine/core';
import React from 'react';
import styled from 'styled-components';
import {ArrowBackIosIcon} from '@mattermost/compass-icons/components';
import {useHistory} from 'react-router-dom';

const HeaderContainer = styled.div`
    display: flex;
    align-items: center;
    justify-content: flex-start;
    width: 100%;
    gap: 8px;
`;

type Params = {
    name: string;
}
const RhsHeader = ({
    name,
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
                    fontSize: '1em',
                    whiteSpace: 'nowrap',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                }}
            >{name || 'Customer Info'}</Title>
        </HeaderContainer>
    );
};

export {
    RhsHeader,
};
