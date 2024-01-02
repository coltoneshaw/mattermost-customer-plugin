import React from 'react';
import {MantineProvider, MantineThemeOverride} from '@mantine/core';

import {MemoryRouter, Route} from 'react-router-dom';

import {CustomerList} from './components/CustomerList/CustomerList';
import {CustomerInfo} from './components/CustomerInfo/CustomerInfo';

const theme: MantineThemeOverride = {
    fontFamily: '"Open Sans", sans-serif;',
    fontSizes: {
        xs: '10px',
        sm: '12px',
        md: '14px',
        lg: '16px',
        xl: '18px',
    },
};

const CustomerRHS = () => {
    const paths = [
        {pathname: '/customers'},
    ];
    return (

        <MemoryRouter
            initialEntries={paths}
            initialIndex={0}
        >
            <MantineProvider
                withNormalizeCSS={true}
                theme={theme}
            >
                <Route path='/customers/:id'>
                    <CustomerInfo/>
                </Route>
                <Route
                    exact={true}
                    path='/customers'
                >
                    <CustomerList/>
                </Route>

            </MantineProvider>
        </MemoryRouter>

    );
};

export {
    CustomerRHS,
};
