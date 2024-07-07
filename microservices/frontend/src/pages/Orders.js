import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './tables.css';

const Orders = () => {
  const [orders, setOrders] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('/api/orders');
        console.log('Fetched orders:', response.data);
        setOrders(response.data);
      } catch (error) {
        console.error('Error fetching orders:', error);
      }
    };

    fetchData();
  }, []);

  return (
    <div>
      <h1>Orders</h1>
      <table className="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>User Name</th>
            <th>User Email</th>
            <th>Product Name</th>
            <th>Product Price</th>
            <th>Product Count</th>
          </tr>
        </thead>
        <tbody>
        {orders.length === 0 ? (
            <tr>
              <td colSpan="3">No data available</td>
            </tr>
          ) : (
          orders.map(order => (
            <tr key={order.order_id}>
              <td>{order.order_id}</td>
              <td>{order.user_name}</td>
              <td>{order.user_email}</td>
              <td>{order.product_name}</td>
              <td>${order.product_price}</td>
              <td>{order.product_count}</td>
            </tr>
          
          )))}
        </tbody>
      </table>
    </div>
  );
};

export default Orders;
