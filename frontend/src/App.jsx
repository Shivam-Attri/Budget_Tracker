import React, { useState, useEffect, createContext, useContext, useMemo, useCallback } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { ChevronLeft, ChevronRight, MoreVertical, PlusCircle, LogOut, LayoutDashboard, ArrowRightLeft, Folder, Target, Settings, Trash2, Edit, X, Search, Calendar as CalendarIcon, Filter as FilterIcon } from 'lucide-react';

// --- API SETUP ---
// This will be populated by Vite from your .env files.
// This approach safely handles environments where import.meta might not be available.
const env = (typeof import.meta !== 'undefined' && import.meta.env) ? import.meta.env : {};
const API_URL = env.VITE_API_URL || 'http://localhost:8000';

// --- CONTEXT PROVIDERS ---
const AuthContext = createContext();
const AppContext = createContext();

const AuthProvider = ({ children }) => {
    const [auth, setAuth] = useState({
        token: localStorage.getItem('accessToken'),
        refreshToken: localStorage.getItem('refreshToken'),
        isAuthenticated: !!localStorage.getItem('accessToken'),
    });

    const login = (accessToken, refreshToken) => {
        localStorage.setItem('accessToken', accessToken);
        localStorage.setItem('refreshToken', refreshToken);
        setAuth({ token: accessToken, refreshToken, isAuthenticated: true });
    };

    const logout = useCallback(async () => {
        const refreshToken = localStorage.getItem('refreshToken');
        if (refreshToken) {
            try {
                await fetch(`${API_URL}/logout`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ refresh_token: refreshToken }),
                });
            } catch (error) {
                console.error("Logout API call failed, proceeding with local logout.", error);
            }
        }
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        setAuth({ token: null, refreshToken: null, isAuthenticated: false });
    }, []);

    const value = useMemo(() => ({
        ...auth,
        login,
        logout,
    }), [auth, logout]);

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

const AppProvider = ({ children }) => {
    const [isTransactionModalOpen, setTransactionModalOpen] = useState(false);
    const [editingTransaction, setEditingTransaction] = useState(null);

    const openTransactionModal = useCallback((transaction = null) => {
        setEditingTransaction(transaction);
        setTransactionModalOpen(true);
    }, []);

    const closeTransactionModal = useCallback(() => {
        setEditingTransaction(null);
        setTransactionModalOpen(false);
    }, []);

    const value = useMemo(() => ({
        isTransactionModalOpen,
        editingTransaction,
        openTransactionModal,
        closeTransactionModal
    }), [isTransactionModalOpen, editingTransaction, openTransactionModal, closeTransactionModal]);

    return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
};

const useAuth = () => useContext(AuthContext);
const useAppContext = () => useContext(AppContext);

// --- API SERVICE WRAPPER (AXIOS-LIKE) ---
const useApi = () => {
    const { token, refreshToken, logout, login } = useAuth();

    const api = useMemo(() => {
        const customFetch = async (url, options = {}) => {
            const getFreshToken = async () => {
                if (!refreshToken) {
                    logout();
                    return null;
                }
                try {
                    const res = await fetch(`${API_URL}/refresh`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ refresh_token: refreshToken }),
                    });
                    if (!res.ok) throw new Error('Refresh failed');
                    const data = await res.json();
                    login(data.access_token, data.refresh_token);
                    return data.access_token;
                } catch (error) {
                    logout();
                    return null;
                }
            };

            const makeRequest = async (accessToken) => {
                const headers = {
                    'Content-Type': 'application/json',
                    ...options.headers,
                };
                if (accessToken) {
                    headers.Authorization = `Bearer ${accessToken}`;
                }
                const res = await fetch(`${API_URL}${url}`, { ...options, headers });

                if (res.status === 401 && refreshToken) {
                    const newAccessToken = await getFreshToken();
                    if (newAccessToken) {
                        return makeRequest(newAccessToken);
                    } else {
                        throw new Error('Session expired');
                    }
                }
                return res;
            };

            return makeRequest(token);
        };
        return {
            get: (url, options) => customFetch(url, { ...options, method: 'GET' }),
            post: (url, body, options) => customFetch(url, { ...options, method: 'POST', body: JSON.stringify(body) }),
            put: (url, body, options) => customFetch(url, { ...options, method: 'PUT', body: JSON.stringify(body) }),
            delete: (url, options) => customFetch(url, { ...options, method: 'DELETE' }),
        };
    }, [token, refreshToken, login, logout]);

    return api;
};


// --- UI COMPONENTS ---

const Spinner = () => (
    <div className="flex justify-center items-center h-full p-8">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500"></div>
    </div>
);

const Card = ({ children, className = '' }) => (
    <div className={`bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 ${className}`}>
        {children}
    </div>
);

const Button = ({ children, onClick, className = '', type = 'button', variant = 'primary', disabled = false }) => {
    const baseClasses = "px-4 py-2 rounded-lg font-semibold focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed";
    const variants = {
        primary: 'bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-500',
        secondary: 'bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-600 focus:ring-gray-500',
        danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500',
    };
    return (
        <button type={type} onClick={onClick} disabled={disabled} className={`${baseClasses} ${variants[variant]} ${className}`}>
            {children}
        </button>
    );
};

const Input = ({ id, type = 'text', placeholder, value, onChange, className = '', required = false }) => (
    <input
        id={id}
        type={type}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        required={required}
        className={`w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-indigo-500 ${className}`}
    />
);

const Select = ({ id, value, onChange, children, className = '', required = false }) => (
    <select
        id={id}
        value={value}
        onChange={onChange}
        required={required}
        className={`w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-indigo-500 ${className}`}
    >
        {children}
    </select>
);

const Modal = ({ isOpen, onClose, title, children }) => {
    if (!isOpen) return null;
    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex justify-center items-center p-4 animate-fade-in-fast" onClick={onClose}>
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-md mx-auto" onClick={(e) => e.stopPropagation()}>
                <div className="flex justify-between items-center p-4 border-b border-gray-200 dark:border-gray-700">
                    <h3 className="text-xl font-bold text-gray-900 dark:text-white">{title}</h3>
                    <button onClick={onClose} className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200">
                        <X size={24} />
                    </button>
                </div>
                <div className="p-6">
                    {children}
                </div>
            </div>
        </div>
    );
};

const ConfirmationDialog = ({ isOpen, onClose, onConfirm, title, message }) => (
    <Modal isOpen={isOpen} onClose={onClose} title={title}>
        <p className="text-gray-600 dark:text-gray-300 mb-6">{message}</p>
        <div className="flex justify-end space-x-4">
            <Button variant="secondary" onClick={onClose}>Cancel</Button>
            <Button variant="danger" onClick={onConfirm}>Confirm</Button>
        </div>
    </Modal>
);

const Sidebar = ({ currentPage, setPage }) => {
    const navItems = [
        { name: 'Dashboard', icon: LayoutDashboard, page: 'Dashboard' },
        { name: 'Transactions', icon: ArrowRightLeft, page: 'Transactions' },
        { name: 'Categories', icon: Folder, page: 'Categories' },
        { name: 'Budgets', icon: Target, page: 'Budgets' },
    ];

    return (
        <aside className="w-64 bg-white dark:bg-gray-800/50 backdrop-blur-sm border-r border-gray-200 dark:border-gray-700 flex-shrink-0 hidden md:block">
            <div className="h-16 flex items-center justify-center text-2xl font-bold text-indigo-600 dark:text-indigo-400">
                Budgetly
            </div>
            <nav className="mt-8 px-4">
                <ul>
                    {navItems.map(item => (
                        <li key={item.name}>
                            <a
                                href="#"
                                onClick={(e) => { e.preventDefault(); setPage(item.page); }}
                                className={`flex items-center px-4 py-3 my-1 rounded-lg transition-colors duration-200 ${currentPage === item.page ? 'bg-indigo-100 dark:bg-indigo-900/50 text-indigo-700 dark:text-indigo-300' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'}`}
                            >
                                <item.icon className="mr-3" size={20} />
                                <span className="font-medium">{item.name}</span>
                            </a>
                        </li>
                    ))}
                </ul>
            </nav>
        </aside>
    );
};

const Header = () => {
    const { logout } = useAuth();
    const { openTransactionModal } = useAppContext();
    return (
        <header className="h-16 bg-white/50 dark:bg-gray-900/50 backdrop-blur-sm border-b border-gray-200 dark:border-gray-700 flex items-center justify-between px-6">
            <div className="md:hidden text-2xl font-bold text-indigo-600 dark:text-indigo-400">Budgetly</div>
            <div className="flex-1"></div>
            <div className="flex items-center space-x-4">
                <Button onClick={() => openTransactionModal()} className="flex items-center">
                    <PlusCircle size={20} className="mr-2" /> New Transaction
                </Button>
                <button onClick={logout} className="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700">
                    <LogOut size={20} className="text-gray-600 dark:text-gray-400" />
                </button>
            </div>
        </header>
    );
};

// --- PAGES ---

const DashboardPage = () => {
    const [summary, setSummary] = useState(null);
    const [loading, setLoading] = useState(true);
    const api = useApi();

    useEffect(() => {
        const fetchSummary = async () => {
            try {
                const res = await api.get('/api/v1/reports/summary');
                if (res.ok) {
                    const data = await res.json();
                    setSummary(data);
                }
            } catch (error) {
                console.error("Failed to fetch summary", error);
            } finally {
                setLoading(false);
            }
        };
        fetchSummary();
    }, [api]);

    if (loading) return <Spinner />;

    const pieData = [
        { name: 'Income', value: summary?.total_income || 0 },
        { name: 'Expense', value: summary?.total_expense || 0 },
    ];
    const COLORS = ['#4f46e5', '#ef4444'];

    return (
        <div className="space-y-6">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <Card className="bg-gradient-to-br from-green-400 to-green-600 text-white">
                    <h3 className="text-lg font-semibold opacity-80">Total Income</h3>
                    <p className="text-4xl font-bold mt-2">${(summary?.total_income || 0).toFixed(2)}</p>
                </Card>
                <Card className="bg-gradient-to-br from-red-400 to-red-600 text-white">
                    <h3 className="text-lg font-semibold opacity-80">Total Expense</h3>
                    <p className="text-4xl font-bold mt-2">${(summary?.total_expense || 0).toFixed(2)}</p>
                </Card>
                <Card className="bg-gradient-to-br from-indigo-400 to-indigo-600 text-white">
                    <h3 className="text-lg font-semibold opacity-80">Net Savings</h3>
                    <p className="text-4xl font-bold mt-2">${(summary?.net_savings || 0).toFixed(2)}</p>
                </Card>
            </div>
            <Card>
                <h3 className="text-xl font-bold mb-4 text-gray-900 dark:text-white">Income vs Expense</h3>
                <div style={{ width: '100%', height: 300 }}>
                    <ResponsiveContainer>
                        <PieChart>
                            <Pie data={pieData} cx="50%" cy="50%" labelLine={false} outerRadius={100} fill="#8884d8" dataKey="value" nameKey="name">
                                {pieData.map((entry, index) => (
                                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                ))}
                            </Pie>
                            <Tooltip formatter={(value) => `$${Number(value).toFixed(2)}`} />
                            <Legend />
                        </PieChart>
                    </ResponsiveContainer>
                </div>
            </Card>
        </div>
    );
};

const TransactionsPage = ({ onDataChange }) => {
    const [transactions, setTransactions] = useState([]);
    const [pagination, setPagination] = useState({});
    const [loading, setLoading] = useState(true);
    const [filters, setFilters] = useState({ page: 1, limit: 10, sort: '-date' });
    const [showFilters, setShowFilters] = useState(false);
    const [deletingId, setDeletingId] = useState(null);
    const api = useApi();
    const { openTransactionModal } = useAppContext();

    const fetchTransactions = useCallback(async () => {
        setLoading(true);
        const query = new URLSearchParams(filters).toString();
        try {
            const res = await api.get(`/api/v1/transactions?${query}`);
            if (res.ok) {
                const data = await res.json();
                setTransactions(data.data || []);
                setPagination(data.metadata || {});
            }
        } catch (error) {
            console.error("Failed to fetch transactions", error);
        } finally {
            setLoading(false);
        }
    }, [api, filters]);

    useEffect(() => {
        fetchTransactions();
    }, [fetchTransactions, onDataChange]);

    const handleDelete = async () => {
        if (!deletingId) return;
        try {
            const res = await api.delete(`/api/v1/transactions/${deletingId}`);
            if (res.ok) {
                fetchTransactions(); // Refetch
            }
        } catch (error) {
            console.error("Failed to delete transaction", error);
        } finally {
            setDeletingId(null);
        }
    };

    const handleFilterChange = (key, value) => {
        setFilters(prev => ({ ...prev, [key]: value, page: 1 }));
    };

    const handlePageChange = (newPage) => {
        if (newPage >= 1 && newPage <= pagination.last_page) {
            setFilters(prev => ({ ...prev, page: newPage }));
        }
    };

    return (
        <div className="space-y-6">
            <ConfirmationDialog
                isOpen={!!deletingId}
                onClose={() => setDeletingId(null)}
                onConfirm={handleDelete}
                title="Delete Transaction"
                message="Are you sure you want to delete this transaction? This action cannot be undone."
            />
            <div className="flex justify-between items-center">
                <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Transactions</h1>
                <Button onClick={() => setShowFilters(!showFilters)}>
                    <FilterIcon size={20} className="mr-2" /> Filters
                </Button>
            </div>
            
            {showFilters && (
                <Card>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <Select onChange={(e) => handleFilterChange('type', e.target.value)} value={filters.type || ''}>
                            <option value="">All Types</option>
                            <option value="income">Income</option>
                            <option value="expense">Expense</option>
                        </Select>
                        <Input type="date" onChange={(e) => handleFilterChange('start_date', e.target.value ? new Date(e.target.value).toISOString() : '')} />
                        <Input type="date" onChange={(e) => handleFilterChange('end_date', e.target.value ? new Date(e.target.value).toISOString() : '')} />
                    </div>
                </Card>
            )}

            <Card>
                {loading ? <Spinner /> : (
                    <div className="overflow-x-auto">
                        <table className="w-full text-left">
                            <thead className="border-b border-gray-200 dark:border-gray-700">
                                <tr>
                                    <th className="p-4">Date</th>
                                    <th className="p-4">Description</th>
                                    <th className="p-4">Type</th>
                                    <th className="p-4 text-right">Amount</th>
                                    <th className="p-4 text-center">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {transactions.map(t => (
                                    <tr key={t.id} className="border-b border-gray-100 dark:border-gray-800">
                                        <td className="p-4 text-sm text-gray-600 dark:text-gray-400">{new Date(t.date).toLocaleDateString()}</td>
                                        <td className="p-4 font-medium text-gray-900 dark:text-white">{t.description}</td>
                                        <td className="p-4">
                                            <span className={`px-2 py-1 text-xs font-semibold rounded-full ${
                                                t.type === 'income' ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                                            }`}>
                                                {t.type}
                                            </span>
                                        </td>
                                        <td className={`p-4 font-mono text-right ${t.type === 'income' ? 'text-green-600' : 'text-red-600'}`}>
                                            ${t.amount.toFixed(2)}
                                        </td>
                                        <td className="p-4 text-center">
                                            <button onClick={() => openTransactionModal(t)} className="p-1 text-gray-500 hover:text-indigo-600"><Edit size={16} /></button>
                                            <button onClick={() => setDeletingId(t.id)} className="p-1 text-gray-500 hover:text-red-600"><Trash2 size={16} /></button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
                <div className="flex justify-between items-center mt-4">
                    <span className="text-sm text-gray-600 dark:text-gray-400">Total Records: {pagination.total_records || 0}</span>
                    <div className="flex items-center space-x-2">
                        <Button variant="secondary" onClick={() => handlePageChange(filters.page - 1)} disabled={filters.page === 1}>
                            <ChevronLeft size={16} />
                        </Button>
                        <span className="text-sm font-medium">Page {pagination.current_page || 1} of {pagination.last_page || 1}</span>
                        <Button variant="secondary" onClick={() => handlePageChange(filters.page + 1)} disabled={!pagination.last_page || filters.page === pagination.last_page}>
                            <ChevronRight size={16} />
                        </Button>
                    </div>
                </div>
            </Card>
        </div>
    );
};

const CategoriesPage = ({ onDataChange }) => {
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setModalOpen] = useState(false);
    const [editingCategory, setEditingCategory] = useState(null);
    const [deletingCategory, setDeletingCategory] = useState(null);
    const api = useApi();

    const fetchCategories = useCallback(async () => {
        setLoading(true);
        try {
            const res = await api.get('/api/v1/categories?limit=100'); // Assuming max 100 categories
            if (res.ok) {
                const data = await res.json();
                setCategories(data.data || []);
            }
        } catch (error) {
            console.error("Failed to fetch categories", error);
        } finally {
            setLoading(false);
        }
    }, [api]);

    useEffect(() => {
        fetchCategories();
    }, [fetchCategories, onDataChange]);

    const handleOpenModal = (category = null) => {
        setEditingCategory(category);
        setModalOpen(true);
    };

    const handleCloseModal = () => {
        setEditingCategory(null);
        setModalOpen(false);
    };

    const handleSave = async (name) => {
        const payload = { name };
        try {
            const res = editingCategory
                ? await api.put(`/api/v1/categories/${editingCategory.id}`, payload)
                : await api.post('/api/v1/categories', payload);
            if (res.ok) {
                fetchCategories();
                handleCloseModal();
            } else {
                const errorData = await res.json();
                alert(errorData.error || "Failed to save category");
            }
        } catch (error) {
            console.error("Failed to save category", error);
        }
    };

    const handleDelete = async () => {
        if (!deletingCategory) return;
        try {
            const res = await api.delete(`/api/v1/categories/${deletingCategory.id}`);
            if (res.ok) {
                fetchCategories();
            } else {
                const errorData = await res.json();
                alert(errorData.error || "Failed to delete category");
            }
        } catch (error) {
            console.error("Failed to delete category", error);
        } finally {
            setDeletingCategory(null);
        }
    };

    return (
        <div className="space-y-6">
            <ConfirmationDialog
                isOpen={!!deletingCategory}
                onClose={() => setDeletingCategory(null)}
                onConfirm={handleDelete}
                title="Delete Category"
                message={`Are you sure you want to delete "${deletingCategory?.name}"? This might affect existing transactions.`}
            />
            <CategoryModal isOpen={isModalOpen} onClose={handleCloseModal} onSave={handleSave} category={editingCategory} />
            <div className="flex justify-between items-center">
                <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Categories</h1>
                <Button onClick={() => handleOpenModal()}>
                    <PlusCircle size={20} className="mr-2" /> New Category
                </Button>
            </div>
            <Card>
                {loading ? <Spinner /> : (
                    <ul className="space-y-3">
                        {categories.map(cat => (
                            <li key={cat.id} className="flex justify-between items-center p-4 rounded-lg bg-gray-50 dark:bg-gray-700/50">
                                <span className="font-medium">{cat.name}</span>
                                <div className="space-x-2">
                                    <button onClick={() => handleOpenModal(cat)} className="p-1 text-gray-500 hover:text-indigo-600"><Edit size={16} /></button>
                                    <button onClick={() => setDeletingCategory(cat)} className="p-1 text-gray-500 hover:text-red-600"><Trash2 size={16} /></button>
                                </div>
                            </li>
                        ))}
                    </ul>
                )}
            </Card>
        </div>
    );
};

const CategoryModal = ({ isOpen, onClose, onSave, category }) => {
    const [name, setName] = useState('');

    useEffect(() => {
        setName(category ? category.name : '');
    }, [category]);

    const handleSubmit = (e) => {
        e.preventDefault();
        onSave(name);
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title={category ? 'Edit Category' : 'New Category'}>
            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label htmlFor="category-name" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Category Name</label>
                    <Input id="category-name" value={name} onChange={(e) => setName(e.target.value)} required />
                </div>
                <div className="flex justify-end space-x-4">
                    <Button variant="secondary" onClick={onClose}>Cancel</Button>
                    <Button type="submit">Save</Button>
                </div>
            </form>
        </Modal>
    );
};

const BudgetsPage = () => {
    // Placeholder for now
    return (
        <div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Budgets</h1>
            <Card className="mt-6">
                <p>Budget management functionality coming soon.</p>
            </Card>
        </div>
    );
};

const TransactionFormModal = ({ onDataChange }) => {
    const { isTransactionModalOpen, closeTransactionModal, editingTransaction } = useAppContext();
    const [transaction, setTransaction] = useState({
        description: '',
        amount: '',
        date: new Date().toISOString().split('T')[0],
        type: 'expense',
        category_id: ''
    });
    const [categories, setCategories] = useState([]);
    const api = useApi();

    useEffect(() => {
        if (isTransactionModalOpen) {
            const fetchCategories = async () => {
                try {
                    const res = await api.get('/api/v1/categories?limit=100');
                    if(res.ok) {
                        const data = await res.json();
                        setCategories(data.data || []);
                        if (editingTransaction) {
                            setTransaction({
                                ...editingTransaction,
                                date: new Date(editingTransaction.date).toISOString().split('T')[0],
                            });
                        } else {
                            // Set default category if available
                            if (data.data && data.data.length > 0) {
                                setTransaction(prev => ({ ...prev, category_id: data.data[0].id }));
                            }
                        }
                    }
                } catch (error) {
                    console.error("Failed to fetch categories for form", error);
                }
            };
            fetchCategories();
        } else {
            // Reset form when modal is closed
            setTransaction({
                description: '',
                amount: '',
                date: new Date().toISOString().split('T')[0],
                type: 'expense',
                category_id: ''
            });
        }
    }, [isTransactionModalOpen, editingTransaction, api]);

    const handleChange = (e) => {
        const { id, value } = e.target;
        setTransaction(prev => ({ ...prev, [id]: value }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        const payload = {
            ...transaction,
            amount: parseFloat(transaction.amount),
        };
        try {
            const res = editingTransaction
                ? await api.put(`/api/v1/transactions/${editingTransaction.id}`, payload)
                : await api.post('/api/v1/transactions', payload);
            
            if (res.ok) {
                onDataChange(); // Trigger refetch in parent
                closeTransactionModal();
            } else {
                const errorData = await res.json();
                alert(errorData.error || "Failed to save transaction");
            }
        } catch (error) {
            console.error("Failed to save transaction", error);
        }
    };

    return (
        <Modal isOpen={isTransactionModalOpen} onClose={closeTransactionModal} title={editingTransaction ? 'Edit Transaction' : 'New Transaction'}>
            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label htmlFor="description" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Description</label>
                    <Input id="description" value={transaction.description} onChange={handleChange} required />
                </div>
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label htmlFor="amount" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Amount</label>
                        <Input id="amount" type="number" value={transaction.amount} onChange={handleChange} required />
                    </div>
                    <div>
                        <label htmlFor="date" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Date</label>
                        <Input id="date" type="date" value={transaction.date} onChange={handleChange} required />
                    </div>
                </div>
                <div>
                    <label htmlFor="type" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Type</label>
                    <Select id="type" value={transaction.type} onChange={handleChange}>
                        <option value="expense">Expense</option>
                        <option value="income">Income</option>
                    </Select>
                </div>
                <div>
                    <label htmlFor="category_id" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Category</label>
                    <Select id="category_id" value={transaction.category_id} onChange={handleChange} required>
                        <option value="" disabled>Select a category</option>
                        {categories.map(cat => (
                            <option key={cat.id} value={cat.id}>{cat.name}</option>
                        ))}
                    </Select>
                </div>
                <div className="flex justify-end space-x-4 pt-4">
                    <Button variant="secondary" onClick={closeTransactionModal}>Cancel</Button>
                    <Button type="submit">Save Transaction</Button>
                </div>
            </form>
        </Modal>
    );
};

const LoginPage = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const { login } = useAuth();
    const [page, setPage] = useState('login');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        try {
            const res = await fetch(`${API_URL}/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password }),
            });
            if (!res.ok) {
                throw new Error('Login failed. Please check your credentials.');
            }
            const data = await res.json();
            login(data.access_token, data.refresh_token);
        } catch (err) {
            setError(err.message);
        }
    };
    
    if (page === 'register') return <RegisterPage setPage={setPage} />;

    return (
        <div className="flex flex-col justify-center items-center min-h-screen bg-gray-50 dark:bg-gray-900">
            <Card className="w-full max-w-md">
                <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-white mb-2">Welcome Back</h2>
                <p className="text-center text-gray-600 dark:text-gray-400 mb-6">Sign in to continue to Budgetly</p>
                {error && <p className="text-red-500 text-center mb-4">{error}</p>}
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Email Address</label>
                        <Input id="email" type="email" placeholder="you@example.com" value={email} onChange={(e) => setEmail(e.target.value)} />
                    </div>
                    <div>
                        <label htmlFor="password" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Password</label>
                        <Input id="password" type="password" placeholder="••••••••" value={password} onChange={(e) => setPassword(e.target.value)} />
                    </div>
                    <Button type="submit" className="w-full !py-3">Sign In</Button>
                </form>
                <p className="text-center text-sm text-gray-600 dark:text-gray-400 mt-6">
                    Don't have an account? <a href="#" onClick={(e) => {e.preventDefault(); setPage('register')}} className="font-medium text-indigo-600 hover:text-indigo-500">Sign up</a>
                </p>
            </Card>
        </div>
    );
};

const RegisterPage = ({ setPage }) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');
        try {
            const res = await fetch(`${API_URL}/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password }),
            });
            if (!res.ok) {
                const errData = await res.text();
                throw new Error(errData || 'Registration failed.');
            }
            setSuccess('Registration successful! Please log in.');
            setTimeout(() => setPage('login'), 2000);
        } catch (err) {
            setError(err.message);
        }
    };
    
    return (
        <div className="flex flex-col justify-center items-center min-h-screen bg-gray-50 dark:bg-gray-900">
            <Card className="w-full max-w-md">
                <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-white mb-2">Create an Account</h2>
                <p className="text-center text-gray-600 dark:text-gray-400 mb-6">Start managing your finances today.</p>
                {error && <p className="text-red-500 text-center mb-4">{error}</p>}
                {success && <p className="text-green-500 text-center mb-4">{success}</p>}
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label htmlFor="email-reg" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Email Address</label>
                        <Input id="email-reg" type="email" placeholder="you@example.com" value={email} onChange={(e) => setEmail(e.target.value)} />
                    </div>
                    <div>
                        <label htmlFor="password-reg" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Password</label>
                        <Input id="password-reg" type="password" placeholder="Minimum 8 characters" value={password} onChange={(e) => setPassword(e.target.value)} />
                    </div>
                    <Button type="submit" className="w-full !py-3">Sign Up</Button>
                </form>
                <p className="text-center text-sm text-gray-600 dark:text-gray-400 mt-6">
                    Already have an account? <a href="#" onClick={(e) => {e.preventDefault(); setPage('login')}} className="font-medium text-indigo-600 hover:text-indigo-500">Sign in</a>
                </p>
            </Card>
        </div>
    );
}

// --- MAIN APP COMPONENT ---

const AppContent = () => {
    const [page, setPage] = useState('Dashboard');
    // A key to force re-fetching data in child components
    const [dataChangeKey, setDataChangeKey] = useState(0);
    const onDataChange = () => setDataChangeKey(prev => prev + 1);

    const renderPage = () => {
        switch (page) {
            case 'Dashboard':
                return <DashboardPage />;
            case 'Transactions':
                return <TransactionsPage onDataChange={dataChangeKey} />;
            case 'Categories':
                return <CategoriesPage onDataChange={dataChangeKey} />;
            case 'Budgets':
                return <BudgetsPage />; // Still a placeholder, but can be built out now
            default:
                return <DashboardPage />;
        }
    };

    return (
        <div className="flex h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-200">
            <Sidebar currentPage={page} setPage={setPage} />
            <div className="flex-1 flex flex-col overflow-hidden">
                <Header />
                <main className="flex-1 overflow-x-hidden overflow-y-auto bg-gray-50 dark:bg-gray-900 p-6">
                    {renderPage()}
                </main>
                <TransactionFormModal onDataChange={onDataChange} />
            </div>
        </div>
    );
};

export default function App() {
    return (
        <AuthProvider>
            <AppProvider>
                <AppWrapper />
            </AppProvider>
        </AuthProvider>
    );
}

const AppWrapper = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <AppContent /> : <LoginPage />;
};
