import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './tables.css';

const Products = () => {
  const [products, setProducts] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('/api/products');
        console.log('Fetched products:', response.data);
        setProducts(response.data);
      } catch (error) {
        console.error('Error fetching products:', error);
      }
    };

    fetchData();
  }, []);

  return (
    <div>
      <h1>Products</h1>
      <table className="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Price</th>
            <th>Count</th>
          </tr>
        </thead>
        <tbody>
        {products.length === 0 ? (
            <tr>
              <td colSpan="3">No data available</td>
            </tr>
          ) : (
          products.map(product => (
            <tr key={product.id}>
              <td>{product.id}</td>
              <td>{product.name}</td>
              <td>${product.price}</td>
              <td>{product.count}</td>
            </tr>
          ))
          )}
        </tbody>
      </table>
    </div>
  );
};

export default Products;
