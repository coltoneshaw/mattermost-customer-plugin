import React from 'react';
import {MantineProvider, MantineThemeOverride} from '@mantine/core';

import {RighthandSidebar} from './components/CustomerList/CustomerList';

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

const CustyRHS = () => {
    return (
        <MantineProvider
            withNormalizeCSS={true}
            theme={theme}
        >
            <RighthandSidebar/>
        </MantineProvider>
    );
};

export {
    CustyRHS,
};
