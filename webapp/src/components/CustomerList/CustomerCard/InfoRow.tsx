import {Badge, Text} from '@mantine/core';
import React from 'react';
import styled from 'styled-components';

import {Customer} from '@/types/customers';
import {getTimestamp} from '@/time/Timestamp';

const InfoRowContainer = styled.div`
display: flex;
justify-content: space-between;
align-items: center;
flexDirection: row;
`;

type InfoRowParams = Pick<Customer, 'lastUpdated' | 'licenseType'>
const InfoRow = ({
    lastUpdated,
    licenseType,
}: InfoRowParams) => {
    return (
        <InfoRowContainer>
            <Text
                mt={0}
                c='dimmed'
                fz='xs'
            >
                {'Last updated: ' + getTimestamp(lastUpdated) }
            </Text>
            { licenseType !== '' && (
                <Badge
                    variant='filled'
                    size='lg'
                    sx={{
                        color: 'rgba(var(--center-channel-color-rgb),0.72)',
                        backgroundColor: 'rgba(var(--center-channel-color-rgb),0.08)',
                        whiteSpace: 'nowrap',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        letterSpacing: '0.5px',
                    }}
                >
                    {licenseType}
                </Badge>
            ) }
        </InfoRowContainer>
    );
};

export {
    InfoRow,
};
