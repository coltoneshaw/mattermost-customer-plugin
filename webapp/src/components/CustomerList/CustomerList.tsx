import React, {useEffect, useState} from 'react';

import {Stack} from '@mantine/core';

import styled from 'styled-components';

import {Customer, CustomerSortOptions, SortDirection} from '@/types/customers';
import {clientFetchCustomers} from '@/client';

import {CustomerCard} from './CustomerCard/CustomerCard';
import {Header} from './Header/Header';

// const customerURL = getApiUrl() + '/customers/';

const Container = styled.div`
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: flex-start;
    align-items: center;
    flex-direction: column;
    overflow-y: auto;
`;

const RighthandSidebar = () => {
    const [customers, setCustomers] = useState<Customer[]>([]);
    const [sortBy, setSortBy] = useState<CustomerSortOptions>(CustomerSortOptions.Default);
    const [orderBy, setOrderBuy] = useState<SortDirection>(SortDirection.DirectionAsc);
    const [page] = useState<number>(0);
    const [perPage] = useState<number>(25);

    useEffect(() => {
        clientFetchCustomers({
            sort: sortBy,
            order: orderBy,
            page: String(page),
            perPage: String(perPage),
        }).
            then((res) => {
                if (!res) {
                    return;
                }
                setCustomers(res.customers);
            });
    }, [orderBy, page, perPage, sortBy]);

    return (
        <Container>
            <Header
                orderBy={orderBy}
                setOrderBy={setOrderBuy}
                sortBy={sortBy}
                setSortBy={setSortBy}
            />
            <Stack
                justify='flex-start'
                align='stretch'
                spacing='sm'
                style={{
                    width: '100%',
                    height: '100%',
                    overflowY: 'auto',
                    padding: '1em',
                }}
            >
                {
                    customers.length > 0 && customers.
                        map((customer) => {
                            return (
                                <CustomerCard
                                    key={customer.id}
                                    customer={customer}
                                />
                            );
                        })
                }
            </Stack>

        </Container>
    );
};

export {RighthandSidebar};
