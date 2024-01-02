import React, {useEffect, useState} from 'react';

import {Stack} from '@mantine/core';

import styled from 'styled-components';

import {Customer, CustomerSortOptions, SortDirection} from '@/types/customers';
import {clientFetchCustomers} from '@/client';

import {CenteredText} from '../CenteredText';

import {CustomerCard} from './CustomerCard/CustomerCard';
import {Header} from './Header/Header';

const Container = styled.div`
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: flex-start;
    align-items: center;
    flex-direction: column;
    overflow-y: auto;
`;

const CustomerList = () => {
    const [customers, setCustomers] = useState<Customer[]>([]);
    const [sortBy, setSortBy] = useState<CustomerSortOptions>(CustomerSortOptions.Default);
    const [orderBy, setOrderBuy] = useState<SortDirection>(SortDirection.DirectionAsc);
    const [page] = useState<number>(0);
    const [perPage] = useState<number>(25);
    const [searchTerm, setSearchTerm] = useState<string>('');

    useEffect(() => {
        clientFetchCustomers({
            sort: sortBy,
            order: orderBy,
            page: String(page),
            perPage: String(perPage),
            searchTerm,
        }).
            then((res) => {
                if (!res) {
                    return;
                }
                setCustomers(res.customers);
            });
    }, [orderBy, page, perPage, sortBy, searchTerm]);

    return (
        <Container>
            <Header
                orderBy={orderBy}
                setOrderBy={setOrderBuy}
                sortBy={sortBy}
                setSortBy={setSortBy}
                setSearchTerm={setSearchTerm}
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
                    (customers && customers.length > 0) ? customers.
                        map((customer) => {
                            return (
                                <CustomerCard
                                    key={customer.id}
                                    customer={customer}
                                />
                            );
                        }) : <CenteredText message={'No customers found'}/>
                }
            </Stack>

        </Container>
    );
};

export {CustomerList};
