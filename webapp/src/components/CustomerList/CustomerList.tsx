import React, {useEffect, useState} from 'react';

import {Stack} from '@mantine/core';

import styled from 'styled-components';

import {useDispatch} from 'react-redux';

import {getMissingProfilesByIds} from 'mattermost-redux/actions/users';

import {Customer, CustomerSortOptions, SortDirection} from '@/types/customers';
import {clientFetchCustomers} from '@/client';

import {CenteredText} from '../CenteredText';

import {RHSTitle} from '../rhsTitle';

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
    const [page] = useState(0);
    const [perPage] = useState(25);
    const [searchTerm, setSearchTerm] = useState('');

    const [loading, setLoading] = useState(false);

    const dispatch = useDispatch();

    useEffect(() => {
        setLoading(true);
        clientFetchCustomers({
            sort: sortBy,
            order: orderBy,
            page: String(page),
            perPage: String(perPage),
            searchTerm,
        }).
            then((res) => {
                setLoading(false);
                if (!res) {
                    return;
                }
                setCustomers(res.customers);
            });

        return () => {
            setLoading(true);
        };
    }, [orderBy, page, perPage, sortBy, searchTerm]);

    useEffect(() => {
        if (customers && customers.length > 0) {
            const ids: string[] = [];
            customers.forEach(({
                technicalAccountManager,
                customerSuccessManager,
                accountExecutive,

            }) => {
                if (technicalAccountManager) {
                    ids.push(technicalAccountManager);
                }
                if (customerSuccessManager) {
                    ids.push(customerSuccessManager);
                }
                if (accountExecutive) {
                    ids.push(accountExecutive);
                }
            });
            dispatch(getMissingProfilesByIds([...new Set(ids)]));
        }
    }, [dispatch, customers]);

    let AlternateContainer = <CenteredText message={'No customers found'}/>;
    if (loading) {
        AlternateContainer = (
            <CenteredText
                message={'Loading...'}
            />
        );
    }

    return (
        <>
            <RHSTitle>
                {'Customers'}
            </RHSTitle>
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
                            }) : AlternateContainer
                    }
                </Stack>

            </Container>
        </>
    );
};

export {CustomerList};
