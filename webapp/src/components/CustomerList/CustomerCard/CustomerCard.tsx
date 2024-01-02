import React from 'react';
import styled from 'styled-components';

import {Text} from '@mantine/core';

import {Customer} from '@/types/customers';

import {UserList} from './UserList';
import {InfoRow} from './InfoRow';
type CustomerParams = {
    customer: Customer;
}

const CustomerCardContainer = styled.div`
width: 100%;
height: min-content;
border: 1px solid rgba(var(--center-channel-color-rgb), 0.16);
background: var(--center-channel-bg);
padding: 16px 20px 20px;
box-shadow: rgba(0, 0, 0, 0.08) 0px 2px 3px 0px;
display: flex;
flex-direction: column;
justify-content: space-between;
gap: 4px;

&:hover {
  background-color: rgba(var(--center-channel-color-rgb), 0.08);
  cursor: pointer;
}
`;

const CustomerCard = ({
    customer,
}: CustomerParams) => {
    return (
        <CustomerCardContainer>
            <Text
                fw={600}
                fz='md'
                mt={0}
                sx={{
                    whiteSpace: 'nowrap',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    maxWidth: '80%',
                }}
            >
                {customer.name}
            </Text>
            <UserList
                accountExecutive={customer.accountExecutive}
                customerSuccessManager={customer.customerSuccessManager}
                technicalAccountManager={customer.technicalAccountManager}
            />
            <InfoRow
                lastUpdated={customer.lastUpdated}
                type={customer.type}
            />
        </CustomerCardContainer>
    );
};

export {
    CustomerCard,
};
