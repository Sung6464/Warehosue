import React, { useState, useEffect, useCallback } from 'react';

// Reusable Button Component with enhanced styling
const Button = ({ children, onClick, className = '', variant = 'primary', size = 'md', ...props }) => {
  let baseStyle = 'inline-flex items-center justify-center rounded-lg text-sm font-semibold transition-all duration-200 ease-in-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none transform active:scale-98';
  let variantStyle = '';
  let sizeStyle = '';

  switch (variant) {
    case 'primary':
      variantStyle = 'bg-gradient-to-r from-blue-600 to-blue-700 text-white shadow-md hover:from-blue-700 hover:to-blue-800 focus:ring-blue-500';
      break;
    case 'outline':
      variantStyle = 'border border-blue-500 text-blue-600 bg-white hover:bg-blue-50 focus:ring-blue-500';
      break;
    case 'destructive':
      variantStyle = 'bg-gradient-to-r from-red-600 to-red-700 text-white shadow-md hover:from-red-700 hover:to-red-800 focus:ring-red-500';
      break;
    case 'secondary':
      variantStyle = 'bg-gray-200 text-gray-800 hover:bg-gray-300 focus:ring-gray-400';
      break;
    case 'ghost':
      variantStyle = 'text-gray-600 hover:bg-gray-100 focus:ring-gray-400';
      break;
    default:
      variantStyle = 'bg-gradient-to-r from-blue-600 to-blue-700 text-white shadow-md hover:from-blue-700 hover:to-blue-800 focus:ring-blue-500';
  }

  switch (size) {
    case 'sm':
      sizeStyle = 'h-9 px-4 text-xs';
      break;
    case 'lg':
      sizeStyle = 'h-12 px-6 text-base';
      break;
    case 'icon':
      sizeStyle = 'h-10 w-10 p-2'; // Adjusted icon button size and padding
      break;
    case 'md':
    default:
      sizeStyle = 'h-10 px-5 py-2.5';
  }

  return (
    <button onClick={onClick} className={`${baseStyle} ${variantStyle} ${sizeStyle} ${className}`} {...props}>
      {children}
    </button>
  );
};

// Reusable Input Component with improved styling
const Input = ({ label, type = 'text', value, onChange, placeholder, className = '', error, ...props }) => (
  <div className="mb-5">
    {label && <label className="block text-gray-700 text-sm font-medium mb-2">{label}</label>}
    <input
      type={type}
      value={value}
      onChange={onChange}
      placeholder={placeholder}
      className={`block w-full px-4 py-2 text-gray-800 bg-white border border-gray-300 rounded-lg shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 transition duration-150 ease-in-out ${error ? 'border-red-500' : ''} ${className}`}
      {...props}
    />
    {error && <p className="text-red-500 text-xs mt-1">{error}</p>}
  </div>
);

// Reusable Select Component
const Select = ({ label, name, value, onChange, options, className = '', error, ...props }) => (
  <div className="mb-5">
    {label && <label htmlFor={name} className="block text-gray-700 text-sm font-medium mb-2">{label}</label>}
    <select
      id={name}
      name={name}
      value={value}
      onChange={onChange}
      className={`block w-full px-4 py-2 text-gray-800 bg-white border border-gray-300 rounded-lg shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 transition duration-150 ease-in-out ${error ? 'border-red-500' : ''} ${className}`}
      {...props}
    >
      <option value="">Select an option</option>
      {options.map((option) => (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      ))}
    </select>
    {error && <p className="text-red-500 text-xs mt-1">{error}</p>}
  </div>
);


// Modal Component
const Modal = ({ show, onClose, title, children }) => {
  if (!show) return null;

  return (
    <div className="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl shadow-2xl w-full max-w-lg overflow-hidden transform transition-all scale-100 ease-out duration-300">
        <div className="flex justify-between items-center px-6 py-4 border-b border-gray-200 bg-gray-50">
          <h3 className="text-xl font-bold text-gray-800">{title}</h3>
          <Button variant="ghost" size="icon" onClick={onClose} aria-label="Close modal">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-x">
              <path d="M18 6 6 18" />
              <path d="m6 6 12 12" />
            </svg>
          </Button>
        </div>
        <div className="p-6">
          {children}
        </div>
      </div>
    </div>
  );
};


// Icon components (Lucide React equivalents using SVG)
const PlusIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-plus">
    <path d="M12 5v14" />
    <path d="M5 12h14" />
  </svg>
);
const EditIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-edit">
    <path d="M22 22H2v-2a6 6 0 0 1 12 0v2" />
    <path d="M12 6H8a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V12" />
    <path d="M16 2L22 8" />
    <path d="M10 12L22 2" />
  </svg>
);
const Trash2Icon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-trash-2">
    <path d="M3 6h18" />
    <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
    <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
    <line x1="10" x2="10" y1="11" y2="17" />
    <line x1="14" x2="14" y1="11" y2="17" />
  </svg>
);
const CheckIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-check">
    <path d="M20 6 9 17l-5-5"/>
  </svg>
);


// Navigation icons
const DashboardIcon = () => ( // Renamed from HomeIcon for clarity with DashboardPage
  <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-home">
    <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
    <polyline points="9 22 9 12 15 12 15 22" />
  </svg>
);
const PackageIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-package">
    <path d="m7.5 4.27 9 5.14" />
    <path d="M2.5 10.36 12 16l9.5-5.64" />
    <path d="M12 22V16" />
    <path d="m2.5 10.36 9.006 5.637a2 2 0 0 0 1.988 0L21.5 10.36" />
    <path d="M12 2L2.5 7.5l9.5 5.5 9.5-5.5L12 2Z" />
  </svg>
);
const BuildingIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-building">
    <rect width="16" height="20" x="4" y="2" rx="2" ry="2" />
    <path d="M9.5 16h5" />
    <path d="M9.5 12h5" />
    <path d="M9.5 8h5" />
  </svg>
);
const UsersIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-users">
    <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 2 0 0 0-4 4v2" />
    <circle cx="9" cy="7" r="4" />
    <path d="M22 21v-2a4 2 0 0 0-3-3.87C18.73 17.65 17.5 19 16 19" />
    <polyline points="17 11 19 13 22 10" />
  </svg>
);
const BoxIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-box">
    <path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z" />
    <path d="m3.3 7 8.7 5 8.7-5" />
    <path d="M12 22V12" />
  </svg>
);


// The API base URL. This is now an ABSOLUTE URL.
const API_BASE_URL = 'http://localhost:8080/api'; 

// Generic fetch utility to handle API calls
const fetchData = async (url, options = {}) => {
  try {
    const response = await fetch(url, options);
    if (!response.ok) {
      // Attempt to parse JSON error message from backend
      const errorData = await response.json().catch(() => ({ message: `HTTP error! status: ${response.status}` }));
      throw new Error(errorData.error || errorData.message || `Failed to fetch data from ${url}. Status: ${response.status}`);
    }
    // Handle 204 No Content responses (like DELETE)
    if (response.status === 204) {
      return null;
    }
    return response.json();
  } catch (error) {
    console.error('API call error:', error);
    throw error; // Re-throw to be caught by component's error state
  }
};

// Base component for displaying lists with add/edit/delete functionality
// Now uses a modal for forms
const CrudPage = ({ title, fields, apiUrl, initialFormState, idField = 'id', children }) => { // Changed idField default to 'id'
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [form, setForm] = useState(initialFormState);
  const [showModal, setShowModal] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [currentId, setCurrentId] = useState(null);
  const [formErrors, setFormErrors] = useState({});

  // Function to fetch all items for the current entity
  const fetchItems = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchData(apiUrl);
      setItems(data || []); // Ensure data is an array
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [apiUrl]);

  // Fetch items on component mount and when apiUrl changes
  useEffect(() => {
    fetchItems();
  }, [fetchItems]);

  // Handle form input changes
  const handleChange = (e) => {
    const { name, value, type } = e.target;
    // For number inputs, parse to int
    const newValue = type === 'number' ? (value === '' ? '' : parseInt(value)) : value;
    setForm((prev) => ({ ...prev, [name]: newValue }));
    // Clear error for the field being changed
    if (formErrors[name]) {
      setFormErrors((prev) => ({ ...prev, [name]: '' }));
    }
  };

  // Validate form
  const validateForm = () => {
    let errors = {};
    let isValid = true;
    fields.forEach(field => {
      // Basic validation: check if field is empty (unless specified as optional)
      // Check if value is truly "empty" for validation purposes
      const fieldValue = form[field.name];
      const isEmpty = fieldValue === '' || fieldValue === null || fieldValue === undefined || (typeof fieldValue === 'number' && isNaN(fieldValue));

      if (!field.optional && isEmpty) { // 'optional' field can be added to field definitions if needed
        errors[field.name] = `${field.label} is required.`;
        isValid = false;
      }
      // Specific validation for number types if not empty
      if (field.type === 'number' && !isEmpty && isNaN(fieldValue)) {
        errors[field.name] = `${field.label} must be a valid number.`;
        isValid = false;
      }
    });
    setFormErrors(errors);
    return isValid;
  };


  // Handle form submission (add new or update existing)
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!validateForm()) {
      return; // Stop if validation fails
    }

    setLoading(true);
    setError(null); // Clear previous errors
    try {
      // Prepare data, converting number fields if necessary before sending to backend
      const dataToSend = {};
      fields.forEach(field => {
          let val = form[field.name];
          if (field.type === 'number' && typeof val === 'string' && val !== '') {
              val = parseInt(val);
          }
          dataToSend[field.name] = val;
      });

      if (isEditing) {
        await fetchData(`${apiUrl}/${currentId}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(dataToSend),
        });
      } else {
        await fetchData(apiUrl, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(dataToSend),
        });
      }
      setShowModal(false); // Close modal on success
      setForm(initialFormState); // Reset form
      setIsEditing(false);
      setCurrentId(null);
      setFormErrors({}); // Clear form errors
      fetchItems(); // Re-fetch items to update the list
    } catch (err) {
      setError(err.message); // Set error message if API call fails
      setLoading(false); // Stop loading on error
    }
  };

  // Open modal for adding new item
  const handleAddNew = () => {
    setForm(initialFormState);
    setIsEditing(false);
    setCurrentId(null);
    setFormErrors({});
    setShowModal(true);
  };

  // Open modal and populate form for editing
  const handleEdit = (item) => {
    // Ensure numbers are numbers, not strings for form population
    const preparedForm = {};
    fields.forEach(field => {
        preparedForm[field.name] = item[field.name];
        if (field.type === 'number' && typeof item[field.name] === 'string') {
            preparedForm[field.name] = parseInt(item[field.name]);
        }
    });
    setForm(preparedForm);
    setIsEditing(true);
    setCurrentId(item[idField]);
    setFormErrors({});
    setShowModal(true);
  };

  // Handle item deletion
  const handleDelete = async (id) => {
    if (window.confirm(`Are you sure you want to delete this ${title.slice(0, -1)}?`)) {
      setLoading(true);
      setError(null);
      try {
        await fetchData(`${apiUrl}/${id}`, { method: 'DELETE' });
        fetchItems();
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }
  };

  // Display loading, error, or content
  if (loading) return <div className="text-center py-12 text-gray-600">Loading {title}...</div>;
  if (error) return <div className="text-center py-12 text-red-600 font-semibold">Error: {error}</div>;

  return (
    <div className="bg-white rounded-xl shadow-lg p-6 md:p-8">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-3xl font-extrabold text-gray-800">{title} Management</h2>
        <Button onClick={handleAddNew} className="flex items-center">
          <PlusIcon className="mr-2" /> Add New {title.slice(0, -1)}
        </Button>
      </div>

      {/* List of Items */}
      <div className="overflow-x-auto rounded-lg border border-gray-200">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {fields.map((field) => (
                <th key={field.name} className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  {field.label}
                </th>
              ))}
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">ID</th>
              <th className="px-6 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {items.length === 0 ? (
              <tr>
                <td colSpan={fields.length + 2} className="text-center py-8 text-gray-500 text-lg">No {title.toLowerCase()} found.</td>
              </tr>
            ) : (
              items.map((item) => (
                <tr key={item[idField]} className="hover:bg-blue-50 transition-colors duration-150 ease-in-out">
                  {fields.map((field) => (
                    <td key={field.name} className="px-6 py-4 whitespace-nowrap text-sm text-gray-800">
                      {/* Special handling for children (e.g., InventoryPage's mapped fields for display) */}
                      {children ? children(item, field.name) : (item[field.name])}
                    </td>
                  ))}
                  <td className="px-6 py-4 whitespace-nowrap text-xs text-gray-500 font-mono">{item[idField]}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                    <div className="flex items-center justify-center space-x-2">
                      <Button variant="secondary" size="icon" onClick={() => handleEdit(item)} title="Edit">
                        <EditIcon />
                      </Button>
                      <Button variant="destructive" size="icon" onClick={() => handleDelete(item[idField])} title="Delete">
                        <Trash2Icon />
                      </Button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Modal for Add/Edit Form */}
      <Modal show={showModal} onClose={() => setShowModal(false)} title={isEditing ? `Edit ${title.slice(0, -1)}` : `Add New ${title.slice(0, -1)}`}>
        <form onSubmit={handleSubmit}>
          {fields.map((field) => {
            // Check if this field should be rendered as a select
            if (field.options) {
              return (
                <Select
                  key={field.name}
                  label={field.label}
                  name={field.name}
                  value={form[field.name] || ''}
                  onChange={handleChange}
                  options={field.options}
                  error={formErrors[field.name]}
                />
              );
            } else {
              return (
                <Input
                  key={field.name}
                  label={field.label}
                  name={field.name}
                  type={field.type || 'text'}
                  value={form[field.name] || ''}
                  onChange={handleChange}
                  placeholder={`Enter ${field.label.toLowerCase()}`}
                  error={formErrors[field.name]}
                />
              );
            }
          })}
          <div className="flex justify-end space-x-3 mt-6">
            <Button type="button" variant="secondary" onClick={() => setShowModal(false)}>
              Cancel
            </Button>
            <Button type="submit">
              <CheckIcon className="mr-2"/> {isEditing ? `Update ${title.slice(0, -1)}` : `Add ${title.slice(0, -1)}`}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

// Specific Pages using CrudPage component
const CustomerPage = () => (
  <CrudPage
    title="Customers"
    apiUrl={`${API_BASE_URL}/customers`}
    fields={[
      { name: 'firstName', label: 'First Name' },
      { name: 'lastName', label: 'Last Name' },
      { name: 'email', label: 'Email', type: 'email' },
      { name: 'phone', label: 'Phone' },
      { name: 'address', label: 'Address' },
    ]}
    initialFormState={{ firstName: '', lastName: '', email: '', phone: '', address: '' }}
  >
    {/* Custom rendering for table to concatenate first and last name */}
    {(item, fieldName) => {
        if (fieldName === 'firstName') return `${item.firstName} ${item.lastName}`;
        if (fieldName === 'lastName') return null; // Don't render separately
        return item[fieldName];
    }}
  </CrudPage>
);


const WarehousePage = () => (
  <CrudPage
    title="Warehouses"
    apiUrl={`${API_BASE_URL}/warehouses`}
    fields={[
      { name: 'name', label: 'Name' },
      { name: 'location', label: 'Location' },
      { name: 'storage', label: 'Storage (capacity)', type: 'number' },
    ]}
    initialFormState={{ name: '', location: '', storage: '' }}
  />
);

const CommodityPage = () => (
  <CrudPage
    title="Commodities"
    apiUrl={`${API_BASE_URL}/commodities`}
    fields={[
      { name: 'name', label: 'Name' },
      { name: 'amount', label: 'Amount (quantity)', type: 'number' },
    ]}
    initialFormState={{ name: '', amount: '' }}
  />
);

// Inventory Page - Now much cleaner using CrudPage and dynamic options
const InventoryPage = () => {
  const [commodities, setCommodities] = useState([]);
  const [loadingLookups, setLoadingLookups] = useState(true);
  const [lookupError, setLookupError] = useState(null);

  // Fetch commodities for the dropdown
  const fetchLookupData = useCallback(async () => {
    setLoadingLookups(true);
    setLookupError(null);
    try {
      const commData = await fetchData(`${API_BASE_URL}/commodities`);
      setCommodities(commData || []);
    } catch (err) {
      setLookupError(err.message);
    } finally {
      setLoadingLookups(false);
    }
  }, []);

  useEffect(() => {
    fetchLookupData();
  }, [fetchLookupData]);

  // Inventory fields now match the Go Inventory model
  const inventoryFields = [
    { name: 'productId', label: 'Product (Commodity)', options: commodities.map(c => ({ value: c.id, label: c.name })) }, // Use 'id' for value
    { name: 'quantity', label: 'Quantity', type: 'number' },
    { name: 'location', label: 'Location' },
  ];

  // Helper to get entity names for display in inventory table
  const getCommodityName = (productId) => {
    const commodity = commodities.find(c => c.id === productId); // Use 'id' to match
    return commodity ? commodity.name : 'N/A';
  };

  if (loadingLookups) return <div className="text-center py-12 text-gray-600">Loading Inventory Lookup Data...</div>;
  if (lookupError) return <div className="text-center py-12 text-red-600 font-semibold">Error loading lookup data: {lookupError}</div>;


  return (
    <CrudPage
      title="Inventory"
      apiUrl={`${API_BASE_URL}/inventory`}
      fields={inventoryFields}
      initialFormState={{ productId: '', quantity: '', location: '' }} // Match Go model fields
    >
      {/* Custom rendering for table cells to show commodity name instead of ID */}
      {(item, fieldName) => {
        switch (fieldName) {
          case 'productId':
            return getCommodityName(item.productId);
          case 'quantity':
            return item.quantity;
          case 'location':
            return item.location;
          default:
            return item[fieldName];
        }
      }}
    </CrudPage>
  );
};


// Main App Component
const App = () => {
  const [currentPage, setCurrentPage] = useState('dashboard'); // Default to dashboard page

  // Utility to add Tailwind CSS and Inter font to the head (for canvas embedding)
  useEffect(() => {
    const script = document.createElement('script');
    script.src = 'https://cdn.tailwindcss.com';
    document.head.appendChild(script);

    const style = document.createElement('style');
    style.innerHTML = `
      @import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap');
      body {
        font-family: 'Inter', sans-serif;
        background-color: #F8F9FB; /* Light background for the overall app */
      }
      /* Animation for dashboard */
      @keyframes fadeInDown {
        from { opacity: 0; transform: translateY(-20px); }
        to { opacity: 1; transform: translateY(0); }
      }
      @keyframes fadeInUp {
        from { opacity: 0; transform: translateY(20px); }
        to { opacity: 1; transform: translateY(0); }
      }
      @keyframes popIn {
        from { opacity: 0; transform: scale(0.9); }
        to { opacity: 1; transform: scale(1); }
      }
      .animate-fade-in-down { animation: fadeInDown 0.6s ease-out forwards; }
      .animate-fade-in-up { animation: fadeInUp 0.6s ease-out 0.2s forwards; }
      .animate-pop-in { animation: popIn 0.5s ease-out 0.4s forwards; }
    `;
    document.head.appendChild(style);

    // Cleanup function for useEffect
    return () => {
      document.head.removeChild(script);
      document.head.removeChild(style);
    };
  }, []);

  // Dashboard Page - A simple landing page
  const DashboardPage = () => (
    <div className="bg-white rounded-xl shadow-lg p-8 text-center min-h-[400px] flex flex-col justify-center items-center">
      <h2 className="text-4xl font-extrabold text-blue-700 mb-4 animate-fade-in-down">
        Welcome to Warehouse Management System
      </h2>
      <p className="text-lg text-gray-600 mb-8 max-w-2xl animate-fade-in-up">
        Efficiently manage your customers, warehouses, commodities, and inventory all from one centralized dashboard.
      </p>
      <Button onClick={() => setCurrentPage('inventory')} size="lg" className="animate-pop-in">
        Start Managing Inventory <BoxIcon className="ml-2" />
      </Button>
    </div>
  );

  // Conditional rendering of pages based on current selection
  const renderPage = () => {
    switch (currentPage) {
      case 'dashboard':
        return <DashboardPage />;
      case 'customers':
        return <CustomerPage />;
      case 'warehouses':
        return <WarehousePage />;
      case 'commodities':
        return <CommodityPage />;
      case 'inventory':
        return <InventoryPage />;
      default:
        return <h1 className="text-xl text-center py-8 text-gray-600">Page not found.</h1>;
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col font-inter antialiased">
      {/* Top Navigation Bar */}
      <nav className="bg-white shadow-lg py-4 px-6 md:px-10 flex flex-col md:flex-row items-center justify-between rounded-b-xl mb-6">
        <div className="flex items-center mb-4 md:mb-0">
          <DashboardIcon className="text-blue-600 mr-3" /> {/* Changed to DashboardIcon */}
          <h1 className="text-3xl font-extrabold text-gray-800">WMS Hub</h1>
        </div>
        <div className="flex flex-wrap justify-center md:justify-end gap-3">
          <Button
            variant={currentPage === 'dashboard' ? 'primary' : 'ghost'}
            onClick={() => setCurrentPage('dashboard')}
            className="flex items-center text-base"
          >
            <DashboardIcon className="mr-2 w-5 h-5" /> Dashboard
          </Button>
          <Button
            variant={currentPage === 'customers' ? 'primary' : 'ghost'}
            onClick={() => setCurrentPage('customers')}
            className="flex items-center text-base"
          >
            <UsersIcon className="mr-2 w-5 h-5" /> Customers
          </Button>
          <Button
            variant={currentPage === 'warehouses' ? 'primary' : 'ghost'}
            onClick={() => setCurrentPage('warehouses')}
            className="flex items-center text-base"
          >
            <BuildingIcon className="mr-2 w-5 h-5" /> Warehouses
          </Button>
          <Button
            variant={currentPage === 'commodities' ? 'primary' : 'ghost'}
            onClick={() => setCurrentPage('commodities')}
            className="flex items-center text-base"
          >
            <PackageIcon className="mr-2 w-5 h-5" /> Commodities
          </Button>
          <Button
            variant={currentPage === 'inventory' ? 'primary' : 'ghost'}
            onClick={() => setCurrentPage('inventory')}
            className="flex items-center text-base"
          >
            <BoxIcon className="mr-2 w-5 h-5" /> Inventory
          </Button>
        </div>
      </nav>

      {/* Main Content Area */}
      <main className="flex-1 p-4 md:p-8 max-w-7xl mx-auto w-full">
        {renderPage()}
      </main>
    </div>
  );
};

export default App;
