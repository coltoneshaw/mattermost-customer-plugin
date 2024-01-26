import React from 'react';

import {MultiSelect} from '@mantine/core';

import {FormTextInputParams} from './types';
import {formSelectInputStyles} from './styles';

const FormMultiSelect = <T, >({
    getInputProps,
    label,
    placeholder,
    formKey,
    data,
    value,
}: FormTextInputParams<T> & {
    data: { value: string, label: string }[],
    value?: string[],
},
) => {
    const {classes} = formSelectInputStyles();
    return (
        <MultiSelect
            bg={'transparent'}
            clearable={true}
            maxDropdownHeight={600}
            classNames={classes}
            placeholder={placeholder || label}
            label={label}
            data={data}
            {...getInputProps(formKey)}
            value={value || getInputProps(formKey).value}
        />
    );
};

export {
    FormMultiSelect,
};
